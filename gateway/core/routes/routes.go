package routes

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	frontendDir := getEnv("FRONTEND_DIR", "../frontend")

	// ===== STATIC PAGES =====
	// Account pages (login, register, etc.)
	router.Static("/assets", filepath.Join(frontendDir, "assets"))
	router.GET("/account/*path", func(c *gin.Context) {
		p := c.Param("path")
		if p == "/" || p == "" {
			p = "/login"
		}
		file := filepath.Join(frontendDir, "account", strings.TrimPrefix(p, "/")+".html")
		if _, err := os.Stat(file); err != nil {
			file = filepath.Join(frontendDir, "account", "login.html")
		}
		c.File(file)
	})

	// Client pages
	router.GET("/client/*path", func(c *gin.Context) {
		p := c.Param("path")
		if p == "/" || p == "" {
			p = "/dashboard"
		}
		file := filepath.Join(frontendDir, "client", strings.TrimPrefix(p, "/")+".html")
		if _, err := os.Stat(file); err != nil {
			file = filepath.Join(frontendDir, "client", "dashboard.html")
		}
		c.File(file)
	})

	// Ops pages
	router.GET("/ops/*path", func(c *gin.Context) {
		p := c.Param("path")
		if p == "/" || p == "" {
			p = "/dashboard"
		}
		file := filepath.Join(frontendDir, "ops", strings.TrimPrefix(p, "/")+".html")
		if _, err := os.Stat(file); err != nil {
			file = filepath.Join(frontendDir, "ops", "dashboard.html")
		}
		c.File(file)
	})

	// Compliance pages
	router.GET("/compliance/*path", func(c *gin.Context) {
		p := c.Param("path")
		if p == "/" || p == "" {
			p = "/dashboard"
		}
		file := filepath.Join(frontendDir, "compliance", strings.TrimPrefix(p, "/")+".html")
		if _, err := os.Stat(file); err != nil {
			file = filepath.Join(frontendDir, "compliance", "dashboard.html")
		}
		c.File(file)
	})

	// Root redirect
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/account/login")
	})

	// ===== API PROXIES =====
	usersURL := getEnv("USERS_SERVICE_URL", "http://localhost:5002")
	notifURL := getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:5003")
	subsURL := getEnv("SUBSCRIPTION_SERVICE_URL", "http://localhost:5004")
	opsURL := getEnv("OPERATIONAL_SERVICE_URL", "http://localhost:5005")

	// Auth APIs → Users service
	router.Any("/api/auth/*path", reverseProxy(usersURL))

	// Ops notification APIs → Notification service
	router.Any("/api/ops/notifications", reverseProxyRewrite(notifURL, "/api/notifications"))
	router.Any("/api/ops/notifications/*path", reverseProxyRewriteWild(notifURL, "/api/ops/notifications", "/api/notifications"))

	// Ops subscription APIs → Subscription service
	router.Any("/api/ops/subscriptions", reverseProxyRewrite(subsURL, "/api/subscriptions"))
	router.Any("/api/ops/subscriptions/*path", reverseProxyRewriteWild(subsURL, "/api/ops/subscriptions", "/api/subscriptions"))

	// Public notifications for clients
	router.GET("/api/notifications/public", reverseProxyRewrite(notifURL, "/api/notifications/public"))

	// Public subscriptions for clients
	router.GET("/api/subscriptions", reverseProxyRewrite(subsURL, "/api/subscriptions"))

	// Compliance KYC APIs → Operational service
	router.Any("/api/compliance/*path", reverseProxy(opsURL))

	// Client KYC APIs → Users service
	router.Any("/api/kyc/*path", reverseProxy(usersURL))

	log.Println("[GATEWAY] Routes configured")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[PROXY ERROR] %s → %s: %v", r.URL.Path, target, err)
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, `{"error":"Service unavailable"}`)
		}
		c.Request.Host = remote.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func reverseProxyRewrite(target, newPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = newPath
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Host = remote.Host
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[PROXY ERROR] %s → %s%s: %v", c.Request.URL.Path, target, newPath, err)
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, `{"error":"Service unavailable"}`)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func reverseProxyRewriteWild(target, oldPrefix, newPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = strings.Replace(c.Request.URL.Path, oldPrefix, newPrefix, 1)
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Host = remote.Host
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[PROXY ERROR] %s → %s: %v", c.Request.URL.Path, target, err)
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, `{"error":"Service unavailable"}`)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
