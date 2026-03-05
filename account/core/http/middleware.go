package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color" // Pastikan sudah go get
	"github.com/gorilla/csrf"
	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/token"
	"github.com/master-abror/zaframework/core/utils"
	"golang.org/x/time/rate" // Pastikan sudah go get
)

type CurrentUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Photo     *string   `json:"photo"`
	GroupID   string    `json:"group_id"`
	LevelID   string    `json:"level_id"`
	RoleID    string    `json:"role_id"`
	Status    string    `json:"status"`
	Level     string    `json:"level"`
	Role      string    `json:"role"`
	Group     string    `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}

// Standard Middleware Chain
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// 1. CORS Middleware (Production Ready)
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	originsMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		originsMap[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if originsMap[origin] || len(allowedOrigins) == 0 {
				// Jika whitelist kosong, hati-hati (dev mode)
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type contextKey string

const CspNonceKey contextKey = "cspNonce"

func WAFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawQuery := strings.ToLower(r.URL.RawQuery)
		attackPatterns := []string{
			"union select", "information_schema", "--", "<script>",
			"javascript:", "/etc/passwd", "exec(",
		}

		for _, pattern := range attackPatterns {
			if strings.Contains(rawQuery, pattern) {
				http.Error(w, "Blocked by WAF", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// 2. Simple WAF Middleware

// func WAFMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// --- 1. WAF LOGIC (Blokir Pola Jahat) ---
// 		rawQuery := strings.ToLower(r.URL.RawQuery)
// 		attackPatterns := []string{
// 			"union select", "information_schema", "--", "<script>",
// 			"javascript:", "/etc/passwd", "exec(", "alert(",
// 		}

// 		for _, pattern := range attackPatterns {
// 			if strings.Contains(rawQuery, pattern) {
// 				http.Error(w, "Request Blocked by Security Filter", http.StatusForbidden)
// 				return // Stop request di sini
// 			}
// 		}

// 		// --- 2. CSP NONCE GENERATION ---
// 		// Buat 16 byte random
// 		nonceBytes := make([]byte, 16)
// 		if _, err := rand.Read(nonceBytes); err != nil {
// 			// Jika gagal generate random (sangat jarang), fallback error
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}
// 		// Convert ke string base64 agar aman di URL/Header
// 		nonce := base64.StdEncoding.EncodeToString(nonceBytes)

// 		// --- 3. SET CSP HEADER ---
// 		// Perhatikan bagian 'nonce-%s'. Ini kuncinya.
// 		// script-src hanya mengizinkan: 'self', CDN tertentu, dan script dengan nonce yang cocok.
// 		cspHeader := fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%s' cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' cdn.jsdelivr.net fonts.googleapis.com; font-src fonts.gstatic.com cdn.jsdelivr.net;", nonce)

// 		w.Header().Set("Content-Security-Policy", cspHeader)

// 		// Header keamanan tambahan (Best Practice)
// 		w.Header().Set("X-Content-Type-Options", "nosniff")
// 		w.Header().Set("X-Frame-Options", "DENY")
// 		w.Header().Set("X-XSS-Protection", "1; mode=block")

// 		// --- 4. MASUKKAN NONCE KE CONTEXT ---
// 		// Agar handler di layer berikutnya (HTML Template) bisa mengambil nonce ini
// 		ctx := context.WithValue(r.Context(), CspNonceKey, nonce)

// 		// Lanjutkan request dengan context baru
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// 3. Rate Limit Middleware (IP Based)
// Menggunakan map sederhana. Untuk production cluster, gunakan Redis/Nginx.
func RateLimitMiddleware(rps float64, burst int) func(http.Handler) http.Handler {
	// Simple map for IP limiter (InMemory)
	// Warning: Ini reset kalau restart.
	var mu sync.Mutex
	visitors := make(map[string]*rate.Limiter)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			// Handle X-Forwarded-For jika di belakang proxy
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = strings.Split(fwd, ",")[0]
			}

			mu.Lock()
			limiter, exists := visitors[ip]
			if !exists {
				limiter = rate.NewLimiter(rate.Limit(rps), burst)
				visitors[ip] = limiter
			}
			mu.Unlock()

			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// 4. Logger Middleware
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap writer untuk tangkap status code
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		status := lrw.statusCode

		// Coloring
		statusColor := color.New(color.FgGreen).SprintFunc()
		if status >= 400 {
			statusColor = color.New(color.FgRed).SprintFunc()
		}

		// Format: [200] POST /path (12ms) IP
		log.Printf("%s %s %s (%v) %s",
			statusColor(status),
			r.Method,
			r.URL.Path,
			duration,
			r.RemoteAddr,
		)
	})
}

// Helper struct untuk logger
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func CSRFMiddleware(authKey string, isProd bool) func(http.Handler) http.Handler {

	opts := []csrf.Option{
		csrf.Path("/"),
		csrf.FieldName("csrf_token"),
		csrf.CookieName("jakedu_csrf"),
		csrf.Secure(isProd), // True jika Production (HTTPS)
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.HttpOnly(true),
	}

	// Error handling custom jika token salah
	opts = append(opts, csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden: CSRF Token Invalid or Missing", http.StatusForbidden)
	})))

	return csrf.Protect([]byte(authKey), opts...)
}

func AuthMiddleware(next http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userID := session.Get(r, "token")

		if userID == nil || userID == "" {

			fmt.Println("user-id-ada-disini", userID)

			code := r.URL.Query().Get("code")
			if code != "" {

				fmt.Println("code-ada-disini", userID)

				fmt.Println("key1267", code)

				// ctx := r.Context()

				// currentUser, err := Authorize(ctx, code)
				// if err != nil {
				// 	ref, err := utils.Encrypt(utils.GetEnv("base_url") + utils.GetEnv("sso_landing"))
				// 	if err != nil {
				// 		log.Printf("%s", err.Error())
				// 	}

				// 	fmt.Println("test123-disini")

				// 	http.Redirect(w, r, utils.GetEnv("sso_auth")+"?ref="+ref, http.StatusSeeOther)
				// 	return
				// }

				_, err := session.Set(w, r, "token", code)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				// ctx = context.WithValue(ctx, "current_user", currentUser)

				// next.ServeHTTP(w, r.WithContext(ctx))

				// return

				ref := utils.GetEnv("base_url") + utils.GetEnv("sso_landing")

				fmt.Println("test123-disini-ye")

				http.Redirect(w, r, ref, http.StatusSeeOther)
				return

			}

			duration := time.Since(start)
			status := http.StatusUnauthorized

			// Coloring
			statusColor := color.New(color.FgGreen).SprintFunc()
			if status >= 400 {
				statusColor = color.New(color.FgRed).SprintFunc()
			}

			log.Printf("%s %s %s (%v) %s",
				statusColor(status),
				r.Method,
				r.URL.Path,
				duration,
				r.RemoteAddr,
			)

			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" || r.Header.Get("Content-Type") == "application/json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"status": false, "msg": "Sesi berakhir, silakan login kembali."}`))
				return
			}

			ref, err := utils.Encrypt(utils.GetEnv("base_url") + utils.GetEnv("sso_landing"))
			if err != nil {
				log.Printf("%s", err.Error())
			}

			fmt.Println("test123-disini")

			http.Redirect(w, r, utils.GetEnv("sso_auth")+"?ref="+ref, http.StatusSeeOther)
			return
		}

		fmt.Println("pass-disini-kok", userID)

		duration := time.Since(start)
		status := http.StatusUnauthorized

		// Coloring
		statusColor := color.New(color.FgGreen).SprintFunc()
		if status >= 400 {
			statusColor = color.New(color.FgRed).SprintFunc()
		}

		log.Printf("%s %s %s (%v) %s",
			statusColor(status),
			r.Method,
			r.URL.Path,
			duration,
			r.RemoteAddr,
		)

		ctx := r.Context()

		currentUser, err := Authorize(r, userID.(string))
		if err != nil {
			fmt.Println("authorize-err", err)
			ref, err := utils.Encrypt(utils.GetEnv("base_url") + utils.GetEnv("sso_landing"))
			if err != nil {
				log.Printf("%s", err.Error())
			}

			http.Redirect(w, r, utils.GetEnv("sso_auth")+"?ref="+ref, http.StatusSeeOther)
		}

		ctx = context.WithValue(ctx, "current_user", currentUser)

		code := r.URL.Query().Get("code")

		if code != "" {
			fmt.Println("disini-tidak-kosong", code)
			ref := utils.GetEnv("base_url") + utils.GetEnv("sso_landing")
			if err != nil {
				log.Printf("%s", err.Error())
			}

			fmt.Println("disini-tidak-kosong", ref)

			http.Redirect(w, r, ref, http.StatusSeeOther)
			return
		}

		fmt.Println("disini-kosong", code)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func UserAuthorize(next http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		currentUser, err := token.GetUserAuthz(r, "token")
		if err != nil {
			http.Redirect(w, r, utils.GetEnv("base_url")+"/account/login", http.StatusSeeOther)
			return
		}

		duration := time.Since(start)
		status := http.StatusOK

		// Coloring
		statusColor := color.New(color.FgGreen).SprintFunc()
		if status >= 400 {
			statusColor = color.New(color.FgRed).SprintFunc()
		}

		log.Printf("%s %s %s (%v) %s",
			statusColor(status),
			r.Method,
			r.URL.Path,
			duration,
			r.RemoteAddr,
		)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "current_user", currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authorize(r *http.Request, code string) (CurrentUser, error) {
	var err error

	fmt.Println("key123", code)
	ctx := r.Context()

	key, err := utils.Decrypt(code)
	if err != nil {
		fmt.Println("decrypt-err", err)
		return CurrentUser{}, err
	}

	// ✅ AMBIL RAW TOKEN DARI REDIS
	jwt, err := utils.RedisGet(ctx, string(key))

	if err != nil {
		fmt.Println("token-err", err)
		DestroyAuthorize(r)
		return CurrentUser{}, err
	}

	fmt.Println("token-login", jwt)

	// ✅ VALIDASI JWT LANGSUNG (TANPA DEKRIPSI SIGNATURE)
	claims, err := utils.ValidateJWT(jwt)
	if err != nil {
		fmt.Println("jwt-validation-error", err)
		return CurrentUser{}, err
	}

	// ✅ REFRESH TOKEN DAN SIMPAN RAW (TANPA ENKRIPSI)
	newToken, _, err := utils.RefreshJWT(jwt, 15*time.Minute)
	if err != nil {
		return CurrentUser{}, err
	}

	err = utils.RedisSet(ctx, string(key), newToken, 15*time.Minute)
	if err != nil {
		return CurrentUser{}, err
	}

	u := claims.User

	photoStr, _ := u["photo"].(string)

	currentUser := CurrentUser{
		ID:      u["id"].(string),
		Name:    u["name"].(string),
		Email:   u["email"].(string),
		Role:    u["role"].(string),
		RoleID:  u["role_id"].(string),
		Level:   u["level"].(string),
		LevelID: u["level_id"].(string),
		Group:   u["group"].(string),
		GroupID: u["group_id"].(string),
		Status:  u["status"].(string),
		Photo:   &photoStr,
	}

	return currentUser, nil
}

func DestroyAuthorize(r *http.Request) error {

	code := session.Get(r, "token")

	ctx := r.Context()

	var err error

	key, err := utils.Decrypt(code.(string))
	if err != nil {
		fmt.Println("decrypt-err", err)
		return err
	}

	err = utils.RedisDel(ctx, string(key))

	if err != nil {
		fmt.Println("token-err", err)
		return err
	}

	return nil
}

func IsLogout(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userID := session.Get(r, "token")

		if userID != nil && userID != "" {

			duration := time.Since(start)
			status := http.StatusUnauthorized

			// Coloring
			statusColor := color.New(color.FgGreen).SprintFunc()
			if status >= 400 {
				statusColor = color.New(color.FgRed).SprintFunc()
			}

			log.Printf("%s %s %s (%v) %s",
				statusColor(status),
				r.Method,
				r.URL.Path,
				duration,
				r.RemoteAddr,
			)
			ref, err := utils.Encrypt(utils.GetEnv("base_url") + utils.GetEnv("sso_landing"))
			if err != nil {
				log.Printf("%s", err.Error())
			}
			http.Redirect(w, r, ref, http.StatusSeeOther)
			return
		}

		duration := time.Since(start)
		status := http.StatusUnauthorized

		// Coloring
		statusColor := color.New(color.FgGreen).SprintFunc()
		if status >= 400 {
			statusColor = color.New(color.FgRed).SprintFunc()
		}

		log.Printf("%s %s %s (%v) %s",
			statusColor(status),
			r.Method,
			r.URL.Path,
			duration,
			r.RemoteAddr,
		)

		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" || r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"status": false, "msg": "Sesi berakhir, silakan login kembali."}`))
			return
		}

		ref, err := utils.Encrypt(utils.GetEnv("base_url") + utils.GetEnv("sso_landing"))
		if err != nil {
			log.Printf("%s", err.Error())
		}

		http.Redirect(w, r, utils.GetEnv("sso_auth")+"?ref="+ref, http.StatusSeeOther)

	})
}
