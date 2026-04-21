package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification/core/concurrency"
	"notification/core/database"
	ehttp "notification/core/http" // Alias package http engine
	"notification/core/utils"      // Wajib import utils untuk GetEnv
)

// Config struct holding all configuration
type Config struct {
	AppName        string
	Port           string
	Env            string // development, production
	AssetsURL      string
	SsoAuth        string
	AllowedOrigins []string

	// DB Configs
	DBConfig database.Config

	// Concurrency Config
	WorkerConfig concurrency.Config
}

// App adalah struct utama framework
type App struct {
	Config     Config
	Router     *http.ServeMux
	Dispatcher *concurrency.Dispatcher
	DB         *database.DBWrapper

	// Helper Response
	Response *ehttp.ResponseHelper
}

// New creates a new App instance
func New(cfg Config) *App {
	var dbWrapper *database.DBWrapper
	var err error

	// =================================================
	// 1. SETUP DATABASE (WITH TOGGLE)
	// =================================================
	// Cek environment variable 'postgres' di .env
	// Jika "true", wajib connect. Jika "false" atau kosong, bypass.
	if utils.GetEnv("postgres") == "true" {
		log.Println("[CORE] Connecting to Database...")
		dbWrapper, err = database.Connect(cfg.DBConfig)
		if err != nil {
			// Jika diset TRUE tapi gagal connect, aplikasi harus mati (Fatal)
			log.Fatalf("[CORE] Database connect error: %v", err)
		}
		log.Println("[CORE] Database Connected ✅")
	} else {
		// Jika diset FALSE, bypass connection
		log.Println("[CORE] Database is DISABLED in .env (Running without DB) ⚠️")
		// Kita buat wrapper tapi Pool-nya nil (aman karena sudah ada safety check di connection.go)
		dbWrapper = &database.DBWrapper{Pool: nil}
	}

	// =================================================
	// 2. SETUP DISPATCHER (WORKER POOL)
	// =================================================
	dispatcher := concurrency.NewDispatcher(cfg.WorkerConfig)
	// Penting: Pastikan dispatcher dijalankan
	dispatcher.Start()

	// =================================================
	// 3. SETUP RESPONSE HELPER
	// =================================================
	// Kita prioritaskan ambil BaseURL dari Config (AllowedOrigins index 0)
	baseURL := os.Getenv("BASE_URL")

	if len(cfg.AllowedOrigins) > 0 {
		baseURL = cfg.AllowedOrigins[0]
	}

	respHelper := &ehttp.ResponseHelper{
		BaseURL:   baseURL,
		AssetsURL: cfg.AssetsURL,
	}

	// Return struct App
	return &App{
		Config:     cfg,
		Router:     http.NewServeMux(),
		Dispatcher: dispatcher,
		DB:         dbWrapper,
		Response:   respHelper,
	}
}

// RegisterJob mendaftarkan background job handler
func (a *App) RegisterJob(name string, handler concurrency.JobHandler) {
	// Mendaftarkan job ke dispatcher
	a.Dispatcher.RegisterJobHandler(name, handler)
}

// Run starts the server with graceful shutdown
func (a *App) Run() {
	// Setup Global Middleware Chain
	// Urutan: RateLimit -> WAF -> CORS -> Logger -> Router
	var handler http.Handler = a.Router

	handler = ehttp.LoggerMiddleware(handler)
	handler = ehttp.CORSMiddleware(a.Config.AllowedOrigins)(handler)
	handler = ehttp.WAFMiddleware(handler)

	// Limit 50 req/sec, burst 100 (Default aman)
	handler = ehttp.RateLimitMiddleware(50, 100)(handler)

	// Setup Server
	srv := &http.Server{
		Addr:         ":" + a.Config.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel error untuk server
	serverErrors := make(chan error, 1)

	if utils.GetEnv("app_secure") == "true" {
		certFile := "certs/localhost.crt"
		keyFile := "certs/localhost.key"
		go func() {
			log.Printf("🚀 %s running on port %s", a.Config.AppName, a.Config.Port)
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				serverErrors <- err
			}
		}()
	} else {
		go func() {
			log.Printf("🚀 %s running on port %s", a.Config.AppName, a.Config.Port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErrors <- err
			}
		}()
	}

	// Start Server in Goroutine

	// Graceful Shutdown Logic
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("[CORE] Server startup failed: %v", err)
	case sig := <-shutdown:
		log.Printf("[CORE] Signal received: %v", sig)
	}

	log.Println("[CORE] Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Stop HTTP
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[CORE] HTTP Shutdown error: %v", err)
	}

	// 2. Stop Workers
	a.Dispatcher.Stop()

	// 3. Close DB
	// Aman dipanggil walau DB mati karena ada safety check di connection.go
	a.DB.Close()

	log.Println("[CORE] Goodbye!")
}
