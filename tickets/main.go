package main

import (
	"log"
	"net/http"
	"time"

	"github.com/master-abror/zaframework/core"
	"github.com/master-abror/zaframework/core/concurrency"
	"github.com/master-abror/zaframework/core/database"
	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/utils"
)

func main() {
	utils.LoadEnv(".env")

	if err := utils.InitJWTLoadKeys("certs/private.pem", "certs/public.pem"); err != nil {
		log.Printf("[WARN] JWT keys not loaded: %v", err)
	}

	if err := utils.InitRedis(); err != nil {
		log.Printf("[WARNING] Redis init failed: %v", err)
	}

	sameSiteMode := http.SameSiteLaxMode
	switch utils.GetEnv("SESSION_SAME_SITE", "Lax") {
	case "Strict":
		sameSiteMode = http.SameSiteStrictMode
	case "None":
		sameSiteMode = http.SameSiteNoneMode
	}

	session.Init(session.Config{
		Driver:      utils.GetEnv("SESSION_DRIVER", "cookie"),
		SecretKey:   utils.GetEnv("SESSION_KEY", utils.GetEnv("APP_KEY")),
		CookieName:  utils.GetEnv("SESSION_NAME", "za_session"),
		SessionLife: utils.ToInt(utils.GetEnv("SESSION_LIFETIME", "3600")),
		Domain:      utils.GetEnv("SESSION_DOMAIN", ""),
		Path:        utils.GetEnv("SESSION_PATH", "/"),
		Secure:      utils.GetEnv("SESSION_SECURE") == "true",
		HttpOnly:    utils.GetEnv("SESSION_HTTP_ONLY") == "true",
		SameSite:    sameSiteMode,
	})

	port := utils.GetEnv("port", "")
	if port == "" {
		port = utils.GetEnv("PORT", "2005")
	}

	cfg := core.Config{
		AppName:        utils.GetEnv("app_name", "Tickets Service"),
		Port:           port,
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
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifetime: 30 * time.Minute,
			MaxConnIdleTime: 10 * time.Minute,
			LogLevel:        "error",
		},

		WorkerConfig: concurrency.Config{
			HighPriorityWorkers:     4,
			NormalPriorityWorkers:   2,
			LowPriorityWorkers:      1,
			QueueSizePerPriority:    1000,
			MaxConcurrentJobs:       500,
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

	app.Router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	app.Run()
}
