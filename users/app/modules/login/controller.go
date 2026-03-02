package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/token"
	"github.com/master-abror/zaframework/core/utils"
)

type Controller struct {
	Dispatcher *concurrency.Dispatcher
	Response   *ehttp.ResponseHelper
}

func NewController(d *concurrency.Dispatcher, r *ehttp.ResponseHelper) *Controller {
	return &Controller{
		Dispatcher: d,
		Response:   r,
	}
}

// Login serve halaman login
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	// Jika sudah login, redirect ke dashboard
	if user, err := token.GetUserAuthz(r, "token"); err == nil && user != nil {
		http.Redirect(w, r, "/account/wrapper", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"AppName": "ThinkNalyze",
		"Year":    time.Now().Year(),
	}
	c.Response.View(w, r, "public/views/login/page.html", "Login to ThinkNalyze", data)
}

// redirectByRole menentukan URL redirect berdasarkan role_code
func redirectByRole(roleCode string) string {
	upper := strings.ToUpper(strings.TrimSpace(roleCode))
	switch upper {
	case "OPERASIONAL":
		return "/ops/dashboard"
	case "COMPLIANCE":
		return "/compliance/dashboard"
	case "CEO":
		return "/ops/dashboard"
	case "CLIENT":
		return "/client/dashboard"
	case "SUPERADMIN":
		return "/ops/dashboard"
	default:
		return "/client/dashboard"
	}
}

// Auth handle login POST request
func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var login Login
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&login); err != nil {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Invalid JSON payload"})
		return
	}

	email := strings.TrimSpace(login.Email)
	password := strings.TrimSpace(login.Password)
	if email == "" || password == "" {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Email dan Password wajib diisi!"})
		return
	}

	browser, os := utils.GetUserAgent(r)

	payload := map[string]any{
		"login_id":  email,
		"password":  password,
		"ip":        utils.GetClientIP(r),
		"browser":   browser,
		"os":        os,
		"latitude":  "",
		"longitude": "",
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "auth", payload, concurrency.PriorityHigh)
	if err != nil {
		fmt.Printf("[LOGIN ERROR] %v\n", err)
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": err.Error()})
		return
	}

	loginRes := result.(LoginResult)
	if loginRes.Token != "" {
		// ✅ SIMPAN TOKEN KE USERS SERVICE SESSION (internal)
		tokenKey, err := token.SetUserAuthz(w, r, loginRes.Token)
		if err != nil {
			fmt.Printf("[LOGIN ERROR] Token SetUserAuthz failed: %v\n", err)
			c.Response.JSON(w, r, map[string]any{"status": false, "msg": err.Error()})
			return
		}

		// ✅ SET COOKIE "token" UNTUK GATEWAY (bisa dibaca cross-service)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenKey,
			Path:     "/",
			Domain:   "",    // kosong = current host (localhost)
			MaxAge:   86400, // 24 jam
			HttpOnly: true,
			Secure:   false, // false untuk development (http://localhost)
			SameSite: http.SameSiteLaxMode,
		})

		fmt.Printf("[LOGIN] Gateway cookie 'token' set with key: %s\n", tokenKey[:20]+"...")

		// ✅ Save RAW token to Redis (backup untuk Gateway)
		saveSessionToRedis(loginRes.UserData, loginRes.Token)
	}

	roleCode, _ := loginRes.UserData["role_code"].(string)
	redirectURL := redirectByRole(roleCode)

	fmt.Printf("[LOGIN SUCCESS] Email: %s, Role: %s, Redirect: %s\n",
		loginRes.UserData["email"], roleCode, redirectURL)

	c.Response.JSON(w, r, map[string]any{
		"status":       true,
		"msg":          "Otentikasi akun berhasil.",
		"redirect_url": redirectURL,
		"user": map[string]any{
			"id":         loginRes.UserData["id"],
			"name":       loginRes.UserData["name"],
			"email":      loginRes.UserData["email"],
			"role":       loginRes.UserData["role"],
			"role_code":  loginRes.UserData["role_code"],
			"group":      loginRes.UserData["group"],
			"group_code": loginRes.UserData["group_code"],
			"level":      loginRes.UserData["level"],
			"level_code": loginRes.UserData["level_code"],
			"photo":      loginRes.UserData["photo"],
		},
	})
}

// PERBAIKAN: Tambah parameter token
func saveSessionToRedis(userData map[string]any, jwtToken string) {
	redisClient := utils.GetRedisClient()
	if redisClient == nil {
		fmt.Println("[LOGIN] Redis not available, skipping session save")
		return
	}

	userID, _ := userData["id"].(string)
	if userID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Simpan JWT token asli (BUKAN yang di-encrypt)
	sessionKey := "session:" + userID
	err := redisClient.Set(ctx, sessionKey, jwtToken, 24*time.Hour).Err()
	if err != nil {
		fmt.Printf("[LOGIN] Failed to save to Redis: %v\n", err)
		return
	}

	fmt.Printf("[LOGIN SESSION] Saved to Redis - Key: %s\n", sessionKey)
}
