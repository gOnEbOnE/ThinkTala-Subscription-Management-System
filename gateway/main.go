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
	Path        string   `json:"path"`
	Target      string   `json:"target"`
	CORS        bool     `json:"cors"`
	Auth        bool     `json:"auth,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Description string   `json:"description,omitempty"`
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
	configPaths := []string{"routes.local.json", "routes.json", "config/routes.json"}
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

		// Inject user info as headers for downstream services
		r.Header.Set("X-User-Role", user.RoleCode)
		r.Header.Set("X-User-ID", user.ID)
		r.Header.Set("X-User-Email", user.Email)

		h.ServeHTTP(w, r)
	}
}

// ========================================
// Role-list Auth Middleware (reads roles from routes.json)
// ========================================
func withRolesAuth(roles []string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUserFromToken(r)
		if err != nil {
			log.Printf("[AUTH] Token error: %v", err)
			if strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Unauthorized","message":"Please log in"}`))
				return
			}
			http.Redirect(w, r, "/account/login", http.StatusFound)
			return
		}

		// SUPERADMIN and CEO bypass role list check
		if user.RoleCode == "SUPERADMIN" || user.RoleCode == "CEO" {
			r.Header.Set("X-User-Role", user.RoleCode)
			r.Header.Set("X-User-ID", user.ID)
			r.Header.Set("X-User-Email", user.Email)
			h.ServeHTTP(w, r)
			return
		}

		// If no roles defined, any authenticated user may pass
		if len(roles) == 0 {
			r.Header.Set("X-User-Role", user.RoleCode)
			r.Header.Set("X-User-ID", user.ID)
			r.Header.Set("X-User-Email", user.Email)
			h.ServeHTTP(w, r)
			return
		}

		// Check if user role is in the allowed list
		for _, allowed := range roles {
			if strings.EqualFold(user.RoleCode, allowed) {
				r.Header.Set("X-User-Role", user.RoleCode)
				r.Header.Set("X-User-ID", user.ID)
				r.Header.Set("X-User-Email", user.Email)
				h.ServeHTTP(w, r)
				return
			}
		}

		log.Printf("[AUTH] Access denied (route roles): user=%s role=%s path=%s allowed=%v",
			user.Email, user.RoleCode, r.URL.Path, roles)

		if strings.Contains(r.Header.Get("Accept"), "application/json") ||
			r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"Forbidden","message":"You don't have access to this resource"}`))
			return
		}

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
<p>You don't have permission to access this resource.<br>Your role: <strong>%s</strong></p>
<a href="/account/login">Back to Login</a>
</div></body></html>`, user.RoleCode)
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
			// Delete session from Redis if token exists
			if tokenCookie, err := r.Cookie("token"); err == nil && tokenCookie.Value != "" {
				if decryptedKey, err := utils.Decrypt(tokenCookie.Value); err == nil {
					ctx := r.Context()
					_ = utils.RedisDel(ctx, string(decryptedKey))
					log.Printf("[GW] Session deleted from Redis on logout")
				}
			}
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

	// Helper to look up target from routes.json config
	getRouteTarget := func(path string) string {
		for _, route := range config.Routes {
			if route.Path == path {
				return route.Target
			}
		}
		return ""
	}

	// --- SUBSCRIPTION SERVICE (role-protected) ---
	// PBI-32,33,34,35,36: Admin package management — CEO, SUPERADMIN, OPERASIONAL only
	subTarget := getRouteTarget("/api/admin/packages")
	if subTarget == "" {
		subTarget = "http://subscription-service.railway.internal:5004"
	}
	http.HandleFunc("/api/admin/packages", withRoleAuth(
		createProxyHandler(subTarget, true),
	))
	http.HandleFunc("/api/admin/packages/", withRoleAuth(
		createProxyHandler(subTarget, true),
	))
	// PBI-37: Public catalog — CLIENT, OPERASIONAL, CEO, SUPERADMIN
	catalogTarget := getRouteTarget("/api/subscription/catalog")
	if catalogTarget == "" {
		catalogTarget = subTarget
	}
	http.HandleFunc("/api/subscription/catalog", withRoleAuth(
		createProxyHandler(catalogTarget, true),
	))
	log.Printf("[GW] Protected API: /api/admin/packages -> %s (CEO/SUPERADMIN/OPERASIONAL)", subTarget)
	log.Printf("[GW] Protected API: /api/subscription/catalog -> %s (CLIENT+)", catalogTarget)

	// --- KYC Admin API (role-protected) - COMPLIANCE & OPERASIONAL can access ---
	kycTarget := getRouteTarget("/api/admin/kyc")
	if kycTarget == "" {
		kycTarget = "http://users-service.railway.internal:2006"
	}
	http.HandleFunc("/api/admin/kyc", withRoleAuth(
		createProxyHandler(kycTarget, true),
	))
	http.HandleFunc("/api/admin/kyc/", withRoleAuth(
		createProxyHandler(kycTarget, true),
	))
	log.Printf("[GW] Protected API: /api/admin/kyc -> %s (CEO/SUPERADMIN/COMPLIANCE/OPERASIONAL)", kycTarget)

	// --- ORDERS Admin API (role-protected) - OPERASIONAL can access ---
	ordersTarget := getRouteTarget("/api/admin/orders")
	if ordersTarget == "" {
		ordersTarget = "http://subscription-service.railway.internal:5004"
	}
	http.HandleFunc("/api/admin/orders", withRolesAuth([]string{"OPERASIONAL", "SUPERADMIN", "CEO"},
		createProxyHandler(ordersTarget, true),
	))
	http.HandleFunc("/api/admin/orders/", withRolesAuth([]string{"OPERASIONAL", "SUPERADMIN", "CEO"},
		createProxyHandler(ordersTarget, true),
	))
	log.Printf("[GW] Protected API: /api/admin/orders -> %s (CEO/SUPERADMIN/OPERASIONAL)", ordersTarget)

	// --- KYC Client API (role-protected) - CLIENT can access their own KYC ---
	kycClientTarget := getRouteTarget("/api/kyc/")
	if kycClientTarget == "" {
		kycClientTarget = kycTarget
	}
	http.HandleFunc("/api/kyc/", withRoleAuth(
		createProxyHandler(kycClientTarget, true),
	))
	log.Printf("[GW] Protected API: /api/kyc/ -> %s (CLIENT+)", kycClientTarget)

	// --- Generic API proxy from routes.json (with optional role auth) ---
	for _, route := range config.Routes {
		if strings.HasPrefix(route.Path, "/api/") {
			// Skip routes already registered above
			if strings.HasPrefix(route.Path, "/api/admin/") ||
				strings.HasPrefix(route.Path, "/api/subscription/") ||
				strings.HasPrefix(route.Path, "/api/kyc") {
				continue
			}
			routeCopy := route // capture loop variable
			handler := createProxyHandler(routeCopy.Target, routeCopy.CORS)
			if routeCopy.Auth {
				http.HandleFunc(routeCopy.Path, withRolesAuth(routeCopy.Roles, handler))
				log.Printf("[GW] Protected API: %s -> %s (roles: %v)", routeCopy.Path, routeCopy.Target, routeCopy.Roles)
			} else {
				http.HandleFunc(routeCopy.Path, handler)
				log.Printf("[GW] API proxy: %s -> %s", routeCopy.Path, routeCopy.Target)
			}
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
