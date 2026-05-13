package session

import (
	"context"
	"fmt"
	"net/http" // Tambahkan ini untuk http.SameSite
	"os"

	"tickets/core/utils"

	"github.com/gorilla/sessions"

	"github.com/rbcervilla/redisstore/v9"
)

// Konfigurasi Session (Diupdate dengan Field Advanced)
type Config struct {
	Driver      string // redis, file, cookie, stateless, stateless_with_redis
	SecretKey   string
	CookieName  string
	SessionLife int

	// Konfigurasi Tambahan (Advanced Security)
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// Manager struct
type Manager struct {
	Config Config
	Store  sessions.Store
}

var SessionManager *Manager

func Init(cfg Config) {
	var store sessions.Store

	// Validasi: Jika driver butuh Redis tapi Redis mati
	if (cfg.Driver == "redis" || cfg.Driver == "stateless_with_redis") && !utils.IsRedisEnabled() {
		fmt.Println("[SESSION] ⚠️ Warning: Driver is set to Redis but Redis is DISABLED. Fallback to 'cookie'.")
		cfg.Driver = "cookie"
	}

	switch cfg.Driver {
	case "redis":
		// Kita ambil Client yang sudah connect dari utils
		client := utils.GetRedisClient()

		// Setup RedisStore pakai client tersebut
		rs, err := redisstore.NewRedisStore(context.Background(), client)
		if err != nil {
			panic("Gagal init Redis Session: " + err.Error())
		}

		// Set Options Lengkap
		rs.Options(sessions.Options{
			Path:     cfg.Path,
			Domain:   cfg.Domain,
			MaxAge:   cfg.SessionLife,
			Secure:   cfg.Secure,
			HttpOnly: cfg.HttpOnly,
			SameSite: cfg.SameSite,
		})
		store = rs

	case "file":
		os.MkdirAll("./storage/sessions", 0755)
		fs := sessions.NewFilesystemStore("./storage/sessions", []byte(cfg.SecretKey))

		// Set Options Lengkap
		fs.Options.Path = cfg.Path
		fs.Options.Domain = cfg.Domain
		fs.Options.MaxAge = cfg.SessionLife
		fs.Options.Secure = cfg.Secure
		fs.Options.HttpOnly = cfg.HttpOnly
		fs.Options.SameSite = cfg.SameSite

		store = fs

	case "cookie":
		cs := sessions.NewCookieStore([]byte(cfg.SecretKey))

		// Set Options Lengkap
		cs.Options.Path = cfg.Path
		cs.Options.Domain = cfg.Domain
		cs.Options.MaxAge = cfg.SessionLife
		cs.Options.Secure = cfg.Secure
		cs.Options.HttpOnly = cfg.HttpOnly
		cs.Options.SameSite = cfg.SameSite

		store = cs

		// stateless & stateless_with_redis tidak butuh Store (Gorilla)
	}

	SessionManager = &Manager{
		Config: cfg,
		Store:  store,
	}

	fmt.Printf("[SESSION] Driver loaded: %s (Secure: %v) 💾\n", cfg.Driver, cfg.Secure)
}
