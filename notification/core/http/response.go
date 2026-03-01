package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type ResponseHelper struct {
	BaseURL   string
	AssetsURL string
	SsoAuth   string
}

// JSONResponse structure (Helper struct untuk standar response)
type JSONResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data,omitempty"`
	Token  string `json:"token,omitempty"`
	CSRF   string `json:"csrf,omitempty"`
}

// GetCSRFToken helper untuk mengambil token string dari request
func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}

// JSON mengirim response dalam format JSON dengan injeksi CSRF otomatis
func (h *ResponseHelper) JSON(w http.ResponseWriter, r *http.Request, data map[string]any) {
	if data == nil {
		data = make(map[string]any)
	}

	// Otomatis Inject CSRF Token ke response JSON
	// Client (Frontend) bisa ambil nilai ini untuk request selanjutnya
	token := csrf.Token(r)
	data["csrf_token"] = token

	w.Header().Set("Content-Type", "application/json")

	// Kita set juga di Header response supaya client mudah ambil (X-CSRF-Token)
	w.Header().Set("X-CSRF-Token", token)

	json.NewEncoder(w).Encode(data)
}

// View merender HTML template dengan Minify dan CSRF Injection
func (h *ResponseHelper) View(w http.ResponseWriter, r *http.Request, filepath string, title string, data map[string]any) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("tester456", h.AssetsURL)

	// ============================================================
	// DATA INJECTION
	// ============================================================
	// Kita ubah Key menjadi PascalCase (Huruf besar diawal)
	// Supaya di HTML bisa dipanggil dengan {{ .BaseURL }}
	viewData := map[string]any{
		"title":      title,
		"base_url":   h.BaseURL,   // Mengambil dari struct (yang diload dari .env)
		"assets_url": h.AssetsURL, // Mengambil dari struct
		// Kita masukkan data dari controller ke key "Data"
		// Di HTML panggilnya: {{ .Data.NamaKey }}
		"data": data,
	}

	// Jika data dari controller ingin langsung diakses di root (opsional)
	// Kita bisa loop map data dan masukkan ke viewData
	maps.Copy(viewData, data)

	// ============================================================
	// CSRF INJECTION
	// ============================================================
	// {{ .csrf_token }} -> mereturn string token
	viewData["csrf_token"] = csrf.Token(r)

	// {{ .csrf_template }} -> mereturn input hidden HTML lengkap
	viewData["csrf_template"] = csrf.TemplateField(r)

	// ============================================================
	// MINIFY & RENDER
	// ============================================================
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, viewData); err != nil {
		http.Error(w, "Render Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	m := minify.New()
	m.Add("text/html", &html.Minifier{})

	minified, err := m.Bytes("text/html", buf.Bytes())
	if err != nil {
		// Fallback to unminified on error (jika minify gagal, tetap tampilkan page)
		minified = buf.Bytes()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(minified)
}

// func (h *ResponseHelper) View(w http.ResponseWriter, r *http.Request, filepath string, title string, data map[string]any) {
// 	// 1. Ambil Nonce dari context (setelah diproses Middleware)
// 	// Pastikan key-nya sama dengan yang ada di Middleware
// 	nonce, ok := r.Context().Value(CspNonceKey).(string)
// 	if !ok {
// 		// Fallback jika middleware lupa dipasang (untuk keamanan, jangan biarkan kosong)
// 		nonce = "default-fallback-nonce"
// 	}

// 	tmpl, err := template.ParseFiles(filepath)
// 	if err != nil {
// 		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// 2. Siapkan View Data
// 	viewData := map[string]any{
// 		"title":      title,
// 		"base_url":   h.BaseURL,
// 		"assets_url": h.AssetsURL,
// 		"nonce":      nonce, // INJECT DISINI: Otomatis bisa dipanggil dengan {{ .Nonce }}
// 		"data":       data,
// 	}

// 	// Copy data controller ke root
// 	for k, v := range data {
// 		viewData[k] = v
// 	}

// 	// CSRF Injection
// 	viewData["csrf_token"] = csrf.Token(r)
// 	viewData["csrf_template"] = csrf.TemplateField(r)

// 	// 3. Render ke Buffer
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, viewData); err != nil {
// 		http.Error(w, "Render Error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// 4. Minify
// 	m := minify.New()
// 	m.Add("text/html", &html.Minifier{})
// 	minified, err := m.Bytes("text/html", buf.Bytes())
// 	if err != nil {
// 		minified = buf.Bytes()
// 	}

// 	// 5. Final Security Headers
// 	// Kita set ulang CSP disini untuk memastikan Nonce benar-benar sinkron dengan template
// 	cspHeader := fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%s' cdn.jsdelivr.net; style-src 'self' cdn.jsdelivr.net fonts.googleapis.com; font-src fonts.gstatic.com cdn.jsdelivr.net;", nonce)

// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.Header().Set("Content-Security-Policy", cspHeader) // Update CSP dengan Nonce terbaru

// 	w.Write(minified)
// }
