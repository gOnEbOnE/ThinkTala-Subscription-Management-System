package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/master-abror/zaframework/app/modules/kyc"
	"github.com/master-abror/zaframework/app/modules/login"
	"github.com/master-abror/zaframework/app/modules/register"
	"github.com/master-abror/zaframework/app/routes"
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

	utils.LoadEnv(".env")

	if err := utils.InitJWTLoadKeys("certs/private.pem", "certs/public.pem"); err != nil {
		log.Fatalf("[FATAL] Gagal memuat kunci JWT: %v", err)
	}

	if err := utils.InitRedis(); err != nil {
		log.Printf("[WARNING] Redis init failed: %v", err)
	}

	// SESSION INIT
	sameSiteMode := http.SameSiteLaxMode
	switch utils.GetEnv("SESSION_SAME_SITE", "Lax") {
	case "Strict":
		sameSiteMode = http.SameSiteStrictMode
	case "None":
		sameSiteMode = http.SameSiteNoneMode
	}

	sessionLife := utils.ToInt(utils.GetEnv("SESSION_LIFETIME", "3600"))

	session.Init(session.Config{
		Driver:      utils.GetEnv("SESSION_DRIVER", "cookie"),
		SecretKey:   utils.GetEnv("SESSION_KEY", utils.GetEnv("APP_KEY")),
		CookieName:  utils.GetEnv("SESSION_NAME", "za_session"),
		SessionLife: sessionLife,

		// Cookie Settings
		Path:     "/",
		Domain:   "", // 👈 PENTING: Set kosong
		Secure:   utils.GetEnv("SESSION_SECURE") == "true",
		HttpOnly: true,
		SameSite: sameSiteMode,
	})

	maxConns, _ := strconv.Atoi(utils.GetEnv("read_db_max_conn", "5"))
	workerMult, _ := strconv.Atoi(utils.GetEnv("APP_WORKER_MULTIPLIER", "4"))

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

	app := core.New(cfg)

	database.MigrateAndSeed(app.DB)

	// ============================================================
	// 2. WIRING FEATURES (Dependency Injection)
	// ============================================================

	// --- Feature: Login ---
	loginRepo := login.NewRepository(app.DB)
	loginService := login.NewService(loginRepo)
	loginController := login.NewController(app.Dispatcher, app.Response)

	// --- Feature: Register ---
	registerRepo := register.NewRepository(app.DB)
	registerService := register.NewService(registerRepo)
	registerController := register.NewController(app.Dispatcher, app.Response)

	// --- Feature: KYC ---
	kycRepo := kyc.NewRepository(app.DB)
	kycService := kyc.NewService(kycRepo)
	kycController := kyc.NewController(app.Dispatcher, app.Response)
	kycAdminController := kyc.NewAdminController(app.Dispatcher, app.Response)

	// Register Job Handlers (Workers)
	app.RegisterJob("auth", loginService.ProcessLoginJob)
	app.RegisterJob("register", registerService.ProcessRegisterJob)
	app.RegisterJob("verify_otp", registerService.ProcessVerifyOTPJob)
	app.RegisterJob("resend_otp", registerService.ProcessResendOTPJob)
	app.RegisterJob("kyc_submit", kycService.ProcessKYCSubmitJob)
	app.RegisterJob("kyc_status", kycService.ProcessKYCStatusJob)
	app.RegisterJob("admin_kyc_list", kycService.ProcessAdminKYCListJob)
	app.RegisterJob("admin_kyc_detail", kycService.ProcessAdminKYCDetailJob)
	app.RegisterJob("admin_kyc_review", kycService.ProcessAdminKYCReviewJob)

	// ============================================================
	// 3. ROUTING
	// ============================================================
	routes.Init(app,
		loginController,
		registerController,
		kycController,
		kycAdminController,
	)

	// ============================================================
	// 4. RUN SERVER
	// ============================================================
	app.Run()
}
