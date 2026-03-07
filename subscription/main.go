package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/master-abror/zaframework/app/modules/packages"
	"github.com/master-abror/zaframework/app/routes"
	"github.com/master-abror/zaframework/core"
	"github.com/master-abror/zaframework/core/concurrency"
	"github.com/master-abror/zaframework/core/database"
	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/utils"
)

func main() {
	// 1. INIT
	utils.LoadEnv(".env")

	if err := utils.InitJWTLoadKeys("certs/private.pem", "certs/public.pem"); err != nil {
		log.Printf("[WARNING] JWT keys not found: %v (continuing without JWT validation)", err)
	}

	if err := utils.InitRedis(); err != nil {
		log.Printf("[WARNING] Redis init failed: %v", err)
	}

	sameSiteMode := http.SameSiteLaxMode
	if utils.GetEnv("SESSION_SAME_SITE", "Lax") == "Strict" {
		sameSiteMode = http.SameSiteStrictMode
	}

	session.Init(session.Config{
		Driver:      utils.GetEnv("SESSION_DRIVER", "cookie"),
		SecretKey:   utils.GetEnv("SESSION_KEY", utils.GetEnv("APP_KEY")),
		CookieName:  utils.GetEnv("SESSION_NAME", "za_session"),
		SessionLife: utils.ToInt(utils.GetEnv("SESSION_LIFETIME", "3600")),
		Path:        "/",
		Domain:      "",
		Secure:      utils.GetEnv("SESSION_SECURE") == "true",
		HttpOnly:    true,
		SameSite:    sameSiteMode,
	})

	maxConns, _ := strconv.Atoi(utils.GetEnv("APP_DB_MAX_CONN", "10"))
	workerMult, _ := strconv.Atoi(utils.GetEnv("APP_WORKER_MULTIPLIER", "4"))

	cfg := core.Config{
		AppName:        utils.GetEnv("app_name", "Thinknalyze Subscription Service"),
		Port:           utils.GetEnv("port", "5004"),
		Env:            utils.GetEnv("app_env", "development"),
		AssetsURL:      utils.GetEnv("assets_url"),
		AllowedOrigins: []string{"*"},

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
			HighPriorityWorkers:   workerMult * 4,
			NormalPriorityWorkers: workerMult * 2,
			LowPriorityWorkers:    workerMult,
			QueueSizePerPriority:  10000,
			MaxConcurrentJobs:     50000,
			MaxRetries:            3,
			JobTimeoutDefault:     10 * time.Second,
			ShutdownTimeout:       30 * time.Second,
			HealthCheckRate:       5 * time.Second,
			MetricsRate:           10 * time.Second,
		},
	}

	app := core.New(cfg)

	// 2. MIGRATION & SEEDING
	database.MigrateAndSeed(app.DB)

	// 3. WIRING
	packagesRepo := packages.NewRepository(app.DB)
	packagesService := packages.NewService(packagesRepo)
	packagesController := packages.NewController(app.Dispatcher, app.Response, packagesService)

	app.RegisterJob("create_package", packagesService.ProcessCreatePackageJob)
	app.RegisterJob("get_admin_packages", packagesService.ProcessGetAdminPackagesJob)
	app.RegisterJob("get_catalog_packages", packagesService.ProcessGetCatalogJob)
	app.RegisterJob("update_package", packagesService.ProcessUpdatePackageJob)
	app.RegisterJob("delete_package", packagesService.ProcessDeletePackageJob)

	// 4. ROUTING
	routes.Init(app, packagesController)

	// 5. RUN
	app.Run()
}
