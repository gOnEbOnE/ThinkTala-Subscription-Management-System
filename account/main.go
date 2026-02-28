package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	// Import Feature Login (Modular)

	// Import Router Baru
	"github.com/master-abror/zaframework/app/modules/account"
	"github.com/master-abror/zaframework/app/modules/dashboard"
	"github.com/master-abror/zaframework/app/routes"

	// Import Framework Core
	"github.com/master-abror/zaframework/core"
	"github.com/master-abror/zaframework/core/concurrency"
	"github.com/master-abror/zaframework/core/database"
	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/utils"
)

func main() {
	// ============================================================
	// 1. INITIALIZATION & CONFIG
	// ============================================================

	// Load Environment Variables
	utils.LoadEnv(".env")

	// Init JWT Keys (Wajib ada file private.pem & public.pem di root)
	if err := utils.InitJWTLoadKeys("private.pem", "public.pem"); err != nil {
		log.Fatalf("[FATAL] Gagal memuat kunci JWT: %v", err)
	}

	// Init Redis (Sesuai request sebelumnya)
	// Jika di .env redis=false, ini akan bypass otomatis tanpa error
	if err := utils.InitRedis(); err != nil {
		log.Printf("[WARNING] Redis init failed: %v", err)
	}

	// ------------------------------------------------------------
	// SESSION INIT (UPDATED)
	// ------------------------------------------------------------
	// 1. Tentukan Mode SameSite (Lax/Strict/None)
	sameSiteMode := http.SameSiteLaxMode // Default
	switch utils.GetEnv("SESSION_SAME_SITE", "Lax") {
	case "Strict":
		sameSiteMode = http.SameSiteStrictMode
	case "None":
		sameSiteMode = http.SameSiteNoneMode
	}

	// 2. Init Session dengan Config Lengkap (Sesuai update manager.go)
	session.Init(session.Config{
		Driver: utils.GetEnv("SESSION_DRIVER", "cookie"),

		// Prioritaskan SESSION_KEY, jika kosong pakai APP_KEY
		SecretKey: utils.GetEnv("SESSION_KEY", utils.GetEnv("APP_KEY")),

		CookieName:  utils.GetEnv("SESSION_NAME", "za_session"),
		SessionLife: utils.ToInt(utils.GetEnv("SESSION_LIFETIME", "3600")),

		// Config Advanced (Dari .env)
		Domain: utils.GetEnv("SESSION_DOMAIN", ""), // Kosongkan jika localhost
		Path:   utils.GetEnv("SESSION_PATH", "/"),

		// Konversi string "true" ke boolean
		Secure:   utils.GetEnv("SESSION_SECURE") == "true",
		HttpOnly: utils.GetEnv("SESSION_HTTP_ONLY") == "true",
		SameSite: sameSiteMode,
	})

	// Helper parsers
	maxConns, _ := strconv.Atoi(utils.GetEnv("read_db_max_conn", "5"))
	workerMult, _ := strconv.Atoi(utils.GetEnv("APP_WORKER_MULTIPLIER", "4"))

	// Config Engine
	cfg := core.Config{
		AppName:        utils.GetEnv("app_name", "ZaFramework"),
		Port:           utils.GetEnv("port", "9002"),
		Env:            utils.GetEnv("app_env", "development"),
		AssetsURL:      utils.GetEnv("assets_url"),
		SsoAuth:        utils.GetEnv("sso_auth"),
		AllowedOrigins: []string{utils.GetEnv("base_url")},

		DBConfig: database.Config{
			Host:            utils.GetEnv("read_db_host"),
			User:            utils.GetEnv("read_db_user"),
			Password:        utils.GetEnv("read_db_pass"),
			DBName:          utils.GetEnv("read_db_name"),
			Port:            utils.GetEnv("read_db_port"),
			SSLMode:         utils.GetEnv("read_db_ssl_mode", "disable"),
			TimeZone:        utils.GetEnv("read_db_timezone", "Asia/Jakarta"),
			MaxConns:        int32(maxConns),
			MinConns:        int32(maxConns / 4),
			MaxConnLifetime: 30 * time.Minute,
			MaxConnIdleTime: 10 * time.Minute,
			LogLevel:        "error",
		},

		WorkerConfig: concurrency.Config{
			HighPriorityWorkers:     workerMult * 4,
			NormalPriorityWorkers:   workerMult * 2,
			LowPriorityWorkers:      workerMult,
			QueueSizePerPriority:    10000,
			MaxConcurrentJobs:       50000,
			QueueFullThreshold:      0.8,
			MaxRetries:              3,
			JobTimeoutDefault:       10 * time.Second,
			ShutdownTimeout:         30 * time.Second,
			EnableAdaptiveRateLimit: true,
			HealthCheckRate:         5 * time.Second,
			MetricsRate:             10 * time.Second,
		},
	}

	// Start Engine
	app := core.New(cfg)

	database.MigrateAndSeed(app.DB)

	// ============================================================
	// 2. WIRING FEATURES (Dependency Injection)
	// ============================================================

	// --- Feature: Login ---
	// accountRepo := account.NewRepository(app.DB)
	// accountService := account.NewService(accountRepo)
	accountController := account.NewController(app.Dispatcher, app.Response)
	dashboardController := dashboard.NewController(app.Dispatcher, app.Response)
	// app.RegisterJob("auth", accountService.ProcessLoginJob)

	// ============================================================
	// 3. ROUTING
	// ============================================================

	// Panggil file router terpisah
	// Kita inject Controller yang sudah di-init di atas ke router
	routes.Init(app,
		accountController,
		dashboardController,
	)

	// ============================================================
	// 4. RUN SERVER
	// ============================================================
	app.Run()
}
