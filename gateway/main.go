package main

import (
	"encoding/json"
	"fmt"
	"gateway/auth"
	"gateway/system"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/master-abror/zaframework/core/utils"
)

func init() {
	tz := system.Env("timezone")
	if tz != "" {
		os.Setenv("TZ", tz)
	}
}

// Route configuration structure
type RouteConfig struct {
	Path        string `json:"path"`
	Target      string `json:"target"`
	CORS        bool   `json:"cors"`
	Description string `json:"description,omitempty"`
}

// Configuration structure
type Config struct {
	AllowedOrigins []string      `json:"allowedOrigins"`
	Routes         []RouteConfig `json:"routes"`
}

var config *Config

// ProxyPool
type ProxyPool struct {
	mu      sync.RWMutex
	proxies map[string]*httputil.ReverseProxy
}

func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		proxies: make(map[string]*httputil.ReverseProxy),
	}
}

func (p *ProxyPool) GetProxy(target string) *httputil.ReverseProxy {
	p.mu.RLock()
	proxy, exists := p.proxies[target]
	p.mu.RUnlock()

	if exists {
		return proxy
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if proxy, exists := p.proxies[target]; exists {
		return proxy
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		log.Printf("Error parsing target URL %s: %v", target, err)
		return nil
	}

	proxy = httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[PROXY ERROR] %s: %v", r.URL.Path, err)
		http.Error(w, "Service temporarily unavailable", http.StatusBadGateway)
	}

	p.proxies[target] = proxy
	return proxy
}

var proxyPool = NewProxyPool()

func loadConfig() (*Config, error) {
	configPaths := []string{"routes.json", "config/routes.json"}
	var configData []byte
	var configPath string

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			configData, _ = ioutil.ReadFile(path)
			if configData != nil {
				break
			}
		}
	}

	if configData == nil {
		return &Config{AllowedOrigins: []string{"*"}}, nil
	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configPath, err)
	}

	log.Printf("Configuration loaded from: %s", configPath)
	log.Printf("Loaded %d routes and %d allowed origins", len(cfg.Routes), len(cfg.AllowedOrigins))
	return &cfg, nil
}

// ========================================
// CORS
// ========================================
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		for _, allowed := range config.AllowedOrigins {
			if strings.TrimSpace(allowed) == "*" || strings.TrimSpace(allowed) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// ========================================
// Logging
// ========================================
func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(ww, r)
		log.Printf("[GW] %s %s %d %v", r.Method, r.URL.Path, ww.statusCode, time.Since(start))
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ========================================
// Role-based Auth Middleware
// ========================================
func withRoleAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AUTH] Path: %s", r.URL.Path)

		user, err := auth.GetUserFromToken(r)
		if err != nil {
			log.Printf("[AUTH] Token error: %v", err)
			http.Redirect(w, r, "/account/login", http.StatusFound)
			return
		}

		log.Printf("[AUTH] User: %s, Role: %s", user.Email, user.RoleCode)

		// Check role access
		if !auth.CheckRoleAccess(r.URL.Path, user.RoleCode) {
			log.Printf("[AUTH] Access denied: user=%s role=%s path=%s", user.Email, user.RoleCode, r.URL.Path)

			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
				strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"Forbidden","message":"You don't have access to this page"}`))
				return
			}

			// Serve 403 page
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<!DOCTYPE html>
<html><head><title>403 Forbidden</title>
<style>
body{font-family:Inter,sans-serif;background:#1a1a2e;color:#fff;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;}
.box{text-align:center;padding:40px;}
h1{font-size:4rem;color:#ff4757;}
a{color:#6C63FF;text-decoration:none;padding:10px 24px;border:2px solid #6C63FF;border-radius:8px;display:inline-block;margin-top:20px;}
a:hover{background:#6C63FF;color:#fff;}
</style></head>
<body><div class="box">
<h1>403</h1>
<h3>Access Denied</h3>
<p>You don't have permission to access this page.<br>Your role: <strong>%s</strong></p>
<a href="/account/login">Back to Login</a>
</div></body></html>`, user.RoleCode)
			return
		}

		log.Printf("[AUTH] Access granted: user=%s role=%s path=%s", user.Email, user.RoleCode, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}

// ========================================
// Proxy handler
// ========================================
func createProxyHandler(target string, enableCORS bool) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		proxy := proxyPool.GetProxy(target)
		if proxy == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		proxy.ServeHTTP(w, r)
	}
	handler = withLogging(handler)
	if enableCORS {
		handler = withCORS(handler)
	}
	return handler
}

// ========================================
// Frontend page server
// ========================================
func serveFrontendPage(frontendDir, section, defaultPage string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := "/" + section + "/"
		pageName := strings.TrimPrefix(r.URL.Path, prefix)
		if pageName == "" {
			pageName = defaultPage
		}
		pageName = strings.TrimSuffix(pageName, "/")

		file := filepath.Join(frontendDir, section, pageName+".html")
		if _, err := os.Stat(file); err == nil {
			http.ServeFile(w, r, file)
			return
		}

		fallback := filepath.Join(frontendDir, section, defaultPage+".html")
		if _, err := os.Stat(fallback); err == nil {
			http.ServeFile(w, r, fallback)
			return
		}

		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

// ========================================
// MAIN
// ========================================
func main() {
	// ✅ HARUS pakai utils.LoadEnv agar utils.GetEnv() bisa kerja
	utils.LoadEnv(".env")

	// Load JWT public key for validation
	if err := utils.InitJWTLoadKeys("certs/private.pem", "certs/public.pem"); err != nil {
		if err2 := utils.InitJWTLoadKeys("../users/certs/private.pem", "../users/certs/public.pem"); err2 != nil {
			log.Fatalf("[FATAL] JWT keys not found: %v / %v", err, err2)
		}
	}

	// ✅ Init Redis SETELAH LoadEnv
	auth.InitRedis()

	// ✅ TAMBAHKAN INI: Load route config
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}

	frontendDir := "../frontend"
	if envDir := system.Env("FRONTEND_DIR"); envDir != "" {
		frontendDir = envDir
	}

	// ========================================
	// 1. Management endpoints (no auth)
	// ========================================
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})
	http.HandleFunc("/admin/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		newConfig, err := loadConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		config = newConfig
		w.Write([]byte(`{"status":"reloaded"}`))
	})

	// ========================================
	// 2. Static assets (no auth)
	// ========================================
	assetsDir := filepath.Join(frontendDir, "assets")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))

	// ========================================
	// 3. Protected Dashboard Pages (with Role Auth)
	// ========================================

	// /ops/* → Only OPERASIONAL, CEO, SUPERADMIN
	http.HandleFunc("/ops/", withRoleAuth(serveFrontendPage(frontendDir, "ops", "dashboard")))

	// /client/* → Only CLIENT, SUPERADMIN
	http.HandleFunc("/client/", withRoleAuth(serveFrontendPage(frontendDir, "client", "dashboard")))

	// /compliance/* → Only COMPLIANCE, SUPERADMIN
	http.HandleFunc("/compliance/", withRoleAuth(serveFrontendPage(frontendDir, "compliance", "dashboard")))

	// ========================================
	// 4. Account pages (no role auth, just serve HTML)
	// ========================================
	http.HandleFunc("/account/", func(w http.ResponseWriter, r *http.Request) {
		// POST /account/login/auth → proxy to users service
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/account/login/auth") {
			for _, route := range config.Routes {
				if route.Path == "/account/login/auth" {
					createProxyHandler(route.Target, route.CORS)(w, r)
					return
				}
			}
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// GET /account/logout
		if r.URL.Path == "/account/logout" {
			// Clear cookies
			http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/", MaxAge: -1})
			http.SetCookie(w, &http.Cookie{Name: "_authz", Value: "", Path: "/", MaxAge: -1})
			http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", Path: "/", MaxAge: -1})
			http.Redirect(w, r, "/account/login", http.StatusFound)
			return
		}

		// Serve HTML pages
		pageName := strings.TrimPrefix(r.URL.Path, "/account/")
		if pageName == "" || pageName == "/" {
			pageName = "login"
		}
		pageName = strings.TrimSuffix(pageName, "/")

		// Try exact match
		file := filepath.Join(frontendDir, "account", pageName+".html")
		if _, err := os.Stat(file); err == nil {
			log.Printf("[GW] Serving: %s -> %s", r.URL.Path, file)
			http.ServeFile(w, r, file)
			return
		}

		// Fallback to login
		fallback := filepath.Join(frontendDir, "account", "login.html")
		log.Printf("[GW] Fallback: %s -> %s", r.URL.Path, fallback)
		http.ServeFile(w, r, fallback)
	})

	// ========================================
	// 5. Root redirect
	// ========================================
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "" {
			// If user has valid token, redirect to their dashboard
			if user, err := auth.GetUserFromToken(r); err == nil {
				redirect := redirectByRole(user.RoleCode)
				http.Redirect(w, r, redirect, http.StatusFound)
				return
			}
			http.Redirect(w, r, "/account/login", http.StatusFound)
			return
		}

		// API routes
		for _, route := range config.Routes {
			if strings.HasPrefix(r.URL.Path, route.Path) {
				createProxyHandler(route.Target, route.CORS)(w, r)
				return
			}
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})

	// ========================================
	// 6. API proxy routes (from routes.json)
	// ========================================

	// --- SUBSCRIPTION SERVICE (role-protected) ---
	// PBI-32,33,34,35,36: Admin package management — CEO, SUPERADMIN, OPERASIONAL only
	http.HandleFunc("/api/admin/packages", withRoleAuth(
		createProxyHandler("http://localhost:5004", true),
	))
	http.HandleFunc("/api/admin/packages/", withRoleAuth(
		createProxyHandler("http://localhost:5004", true),
	))
	// PBI-37: Public catalog — CLIENT, OPERASIONAL, CEO, SUPERADMIN
	http.HandleFunc("/api/subscription/catalog", withRoleAuth(
		createProxyHandler("http://localhost:5004", true),
	))
	log.Printf("[GW] Protected API: /api/admin/packages -> http://localhost:5004 (CEO/SUPERADMIN/OPERASIONAL)")
	log.Printf("[GW] Protected API: /api/subscription/catalog -> http://localhost:5004 (CLIENT+)")

	// --- Generic API proxy from routes.json (no role auth) ---
	for _, route := range config.Routes {
		if strings.HasPrefix(route.Path, "/api/") {
			// Skip routes already registered above
			if strings.HasPrefix(route.Path, "/api/admin/") ||
				strings.HasPrefix(route.Path, "/api/subscription/") {
				continue
			}
			http.HandleFunc(route.Path, createProxyHandler(route.Target, route.CORS))
			log.Printf("[GW] API proxy: %s -> %s", route.Path, route.Target)
		}
	}

	// ========================================
	// 7. Start HTTP server
	// ========================================
	port := system.Env("port")
	if port == "" {
		port = "2000"
	}

	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("=========================================")
	fmt.Printf("  Gateway started on http://localhost:%s\n", port)
	fmt.Println("=========================================")
	fmt.Println("  Dashboard URLs (role-protected):")
	fmt.Printf("    Ops:        http://localhost:%s/ops/dashboard\n", port)
	fmt.Printf("    Client:     http://localhost:%s/client/dashboard\n", port)
	fmt.Printf("    Compliance: http://localhost:%s/compliance/dashboard\n", port)
	fmt.Println("  Login:        http://localhost:" + port + "/account/login")
	fmt.Println("=========================================")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func redirectByRole(roleCode string) string {
	switch strings.ToUpper(roleCode) {
	case "OPERASIONAL", "CEO":
		return "/ops/dashboard"
	case "COMPLIANCE":
		return "/compliance/dashboard"
	case "CLIENT":
		return "/client/dashboard"
	case "SUPERADMIN":
		return "/ops/dashboard"
	default:
		return "/client/dashboard"
	}
}
