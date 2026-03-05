package landing

import (
	"net/http"
	"time"

	// Import package http core kamu
	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
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

func (c *Controller) Welcome(w http.ResponseWriter, r *http.Request) {
	// 1. Siapkan Data
	// Note: Title tidak perlu masuk sini karena ada parameter khusus di method View
	data := map[string]any{
		"AppName": "ZaFramework",
		"Year":    time.Now().Year(),
	}

	// 2. Panggil Response.View (Punya fitur CSRF & Minify)
	// Parameter: (w, r, FilePath, Title, DataMap)
	// Pastikan path file lengkap relatif dari root project
	c.Response.View(w, r, "public/views/landing/index.html", "Welcome to ZAFramework", data)
}

func (c *Controller) MacroData(w http.ResponseWriter, r *http.Request) {
	// 1. Siapkan Data
	// Note: Title tidak perlu masuk sini karena ada parameter khusus di method View
	data := map[string]any{
		"AppName": "ZaFramework",
		"Year":    time.Now().Year(),
	}

	// 2. Panggil Response.View (Punya fitur CSRF & Minify)
	// Parameter: (w, r, FilePath, Title, DataMap)
	// Pastikan path file lengkap relatif dari root project
	c.Response.View(w, r, "public/views/landing/macrodata.html", "Macro Data Dashboard", data)
}

func (c *Controller) ThinkArah(w http.ResponseWriter, r *http.Request) {
	// 1. Siapkan Data
	// Note: Title tidak perlu masuk sini karena ada parameter khusus di method View
	data := map[string]any{
		"AppName": "ZaFramework",
		"Year":    time.Now().Year(),
	}

	// 2. Panggil Response.View (Punya fitur CSRF & Minify)
	// Parameter: (w, r, FilePath, Title, DataMap)
	// Pastikan path file lengkap relatif dari root project
	c.Response.View(w, r, "public/views/landing/thinkarah.html", "ThinkArah", data)
}

func (c *Controller) Docs(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"AppName": "ZaFramework",
	}
	// Render file docs/index.html yang baru kita buat
	c.Response.View(w, r, "public/views/docs/index.html", "Documentation", data)
}
