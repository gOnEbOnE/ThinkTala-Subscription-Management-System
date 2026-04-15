package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	// Import package http core kamu
	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/session"
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

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Siapkan Data
	// Note: Title tidak perlu masuk sini karena ada parameter khusus di method

	data := map[string]any{
		"AppName": "ZaFramework",
		"Year":    time.Now().Year(),
	}

	// 2. Panggil Response.View (Punya fitur CSRF & Minify)
	// Parameter: (w, r, FilePath, Title, DataMap)
	// Pastikan path file lengkap relatif dari root project
	c.Response.View(w, r, "public/views/login/page.html", "Login to ZAFramework", data)
}

func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {

	// Pastikan JSON
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	// Batasi ukuran body
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	var login Login
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // anti field liar

	if err := decoder.Decode(&login); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
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
		"latitude":  "", // Sudah variabel string
		"longitude": "", // Sudah variabel string
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "auth", payload, concurrency.PriorityHigh)

	if err != nil {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": err.Error()})
		return
	}

	loginRes, ok := result.(LoginResult)
	if !ok {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Internal Error casting result"})
		return
	}

	if loginRes.Token != "" {
		_, err := session.Set(w, r, "token", loginRes.Token)
		if err != nil {
			fmt.Println("error", err)
			c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Gagal menyimpan session: "})
			return
		}
	}

	c.Response.JSON(w, r, map[string]any{
		"status": true,
		"msg":    "Otentikasi akun berhasil.",
	})

}
