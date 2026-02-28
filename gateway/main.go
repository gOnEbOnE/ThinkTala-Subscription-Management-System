package main

import (
	"encoding/json"
	"fmt"
	"gateway/system"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	os.Setenv("TZ", system.Env("timezone"))
}

// Route configuration structure
type RouteConfig struct {
	Path        string `json:"path"`
	Target      string `json:"target"`
	CORS        bool   `json:"cors"`
	Description string `json:"description,omitempty"`
}

// Configuration structure untuk JSON file
type Config struct {
	AllowedOrigins []string      `json:"allowedOrigins"`
	Routes         []RouteConfig `json:"routes"`
}

// Global config variable
var config *Config

// ProxyPool untuk reuse proxy instances
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

	// Double-check setelah acquire write lock
	if proxy, exists := p.proxies[target]; exists {
		return proxy
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		log.Printf("Error parsing target URL %s: %v", target, err)
		return nil
	}

	proxy = httputil.NewSingleHostReverseProxy(targetURL)

	// Konfigurasi proxy dengan timeout dan error handling
	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	// Error handler untuk proxy
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for %s: %v", r.URL.Path, err)
		http.Error(w, "Service temporarily unavailable", http.StatusBadGateway)
	}

	p.proxies[target] = proxy
	return proxy
}

var proxyPool = NewProxyPool()

// Load configuration from JSON file
func loadConfig() (*Config, error) {
	// Coba beberapa lokasi file
	configPaths := []string{
		"routes.json",
		"config/routes.json",
		"/etc/jakedu/routes.json",
	}

	var configData []byte

	var configPath string

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			configData, err = ioutil.ReadFile(path)
			if err == nil {
				break
			}
		}
	}

	if configData == nil {
		log.Println("Warning: routes.json not found, using default configuration")
	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file %s: %v", configPath, err)
	}

	log.Printf("Configuration loaded from: %s", configPath)
	log.Printf("Loaded %d routes and %d allowed origins", len(cfg.Routes), len(cfg.AllowedOrigins))

	return &cfg, nil
}

// Reload configuration (untuk hot reload tanpa restart)
func reloadConfig() error {
	newConfig, err := loadConfig()
	if err != nil {
		return err
	}
	config = newConfig
	log.Println("Configuration reloaded successfully")
	return nil
}

// Improved CORS handler dengan security enhancements
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Validasi origin dari config
		originAllowed := false
		for _, allowed := range config.AllowedOrigins {
			if strings.TrimSpace(allowed) == origin {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}

// Middleware untuk logging dan monitoring
func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Custom ResponseWriter untuk capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		h.ServeHTTP(ww, r)

		log.Printf("%s %s %d %v", r.Method, r.URL.Path, ww.statusCode, time.Since(start))
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

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
}

// Config reload endpoint (untuk admin)
func configReload(w http.ResponseWriter, r *http.Request) {
	// Basic auth atau API key validation bisa ditambahkan di sini
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := reloadConfig(); err != nil {
		log.Printf("Failed to reload config: %v", err)
		http.Error(w, "Failed to reload configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Configuration reloaded","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
}

// Routes info endpoint
func routesInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	info := map[string]interface{}{
		"totalRoutes":    len(config.Routes),
		"allowedOrigins": config.AllowedOrigins,
		"routes":         config.Routes,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(info)
}

// Generic proxy handler - Updated to use direct target URLs
func createProxyHandler(target string, enableCORS bool) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if target == "" {
			log.Printf("Target URL not configured")
			http.Error(w, "Service configuration error", http.StatusInternalServerError)
			return
		}

		proxy := proxyPool.GetProxy(target)
		if proxy == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		proxy.ServeHTTP(w, r)
	}

	// Apply middlewares
	handler = withLogging(handler)
	if enableCORS {
		handler = withCORS(handler)
	}

	return handler
}

func main() {
	// Load configuration dari JSON file
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Management endpoints
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/admin/config/reload", configReload)
	http.HandleFunc("/admin/routes/info", routesInfo)
	fs := http.FileServer(http.Dir("./public/jakedu/assets"))
	http.Handle("/j-assets/", http.StripPrefix("/j-assets/", fs))

	// Register all routes dynamically dari JSON config
	for _, route := range config.Routes {
		http.HandleFunc(route.Path, createProxyHandler(route.Target, route.CORS))
		log.Printf("Route registered: %s -> %s (CORS: %t)", route.Path, route.Target, route.CORS)
	}

	// Server configuration
	port := system.Env("port")
	if port == "" {
		port = "2000"
	}

	// Custom logger yang filter TLS handshake errors
	customLogger := log.New(os.Stderr, "", log.LstdFlags)

	// Create HTTP server with timeouts dan custom error log
	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		ErrorLog:       customLogger,
	}

	// Filter TLS handshake errors dari log
	log.SetOutput(&filteredWriter{original: os.Stderr})

	fmt.Printf("Server started at port: %s\n", port)
	fmt.Printf("Health check available at: http://localhost:%s/health\n", port)
	fmt.Printf("Routes info available at: http://localhost:%s/admin/routes/info\n", port)
	fmt.Printf("Config reload available at: http://localhost:%s/admin/config/reload\n", port)
	fmt.Printf("Loaded %d routes from configuration\n", len(config.Routes))

	certFile := "certs/localhost.crt"
	keyFile := "certs/localhost.key"

	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}

// Custom writer untuk filter TLS errors
type filteredWriter struct {
	original *os.File
}

func (fw *filteredWriter) Write(p []byte) (n int, err error) {
	s := string(p)
	// Skip TLS handshake error messages
	if strings.Contains(s, "TLS handshake error") &&
		strings.Contains(s, "client sent an HTTP request to an HTTPS server") {
		return len(p), nil // Don't write to log
	}
	return fw.original.Write(p)
}

/*
	Copyright © 2024 - 2025. PT Arunika Tala Archipelago
	Developed by Muhammad Abror
	Optimized version with direct IP configuration
	For more info, please visit https://arunikatala.co.id
*/
