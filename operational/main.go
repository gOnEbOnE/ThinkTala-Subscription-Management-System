package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	// Import Feature Login (Modular)

	"github.com/master-abror/zaframework/app/modules/account/dashboard"
	"github.com/master-abror/zaframework/app/modules/account/users"
	"github.com/master-abror/zaframework/app/modules/account/wrapper"
	"github.com/master-abror/zaframework/app/modules/landing"
	"github.com/master-abror/zaframework/app/modules/login"
	"github.com/master-abror/zaframework/app/modules/register"
	"github.com/master-abror/zaframework/app/modules/reset"

	// Import Router Baru
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
	maxConns, _ := strconv.Atoi(utils.GetEnv("APP_DB_MAX_CONN", "100"))
	workerMult, _ := strconv.Atoi(utils.GetEnv("APP_WORKER_MULTIPLIER", "4"))

	// Config Engine
	cfg := core.Config{
		AppName:        utils.GetEnv("app_name", "Jakedu Login Service"),
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

	// ============================================================
	// 2. WIRING FEATURES (Dependency Injection)
	// ============================================================

	// --- Feature: Login ---

	// loginCtrl := login.NewController(app.Dispatcher, app.Response)

	landingCtrl := landing.NewController(app.Dispatcher, app.Response)

	loginRepo := login.NewRepository(app.DB)
	loginService := login.NewService(loginRepo)
	loginController := login.NewController(app.Dispatcher, app.Response)

	registerController := register.NewController(app.Dispatcher, app.Response)
	resetController := reset.NewController(app.Dispatcher, app.Response)

	wrapperController := wrapper.NewController(app.Dispatcher, app.Response)
	dashboardController := dashboard.NewController(app.Dispatcher, app.Response)

	usersRepo := users.NewRepository(app.DB)
	usersService := users.NewService(usersRepo)
	usersController := users.NewController(app.Dispatcher, app.Response)

	// Register Job Handler (Worker)
	// Logic background process didaftarkan di sini
	app.RegisterJob("auth", loginService.ProcessLoginJob)
	app.RegisterJob("get_inactive_users", usersService.GetInactiveUsersService)
	// Di file registry worker Anda
	app.RegisterJob("get_detail_user", usersService.ProcessGetDetailUserJob)
	app.RegisterJob("update_user", usersService.ProcessUpdateUserJob)

	// ============================================================
	// 3. ROUTING
	// ============================================================

	// Panggil file router terpisah
	// Kita inject Controller yang sudah di-init di atas ke router
	routes.Init(app,
		landingCtrl,
		loginController,
		registerController,
		resetController,
		wrapperController,
		dashboardController,
		usersController,
	)

	// ============================================================
	// 4. RUN SERVER
	// ============================================================
	app.Run()
}
