package wrapper

import (
	"fmt"
	"net/http"
	"strings"

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

func (c *Controller) Welcome(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/wrapper/page.html", "Welcome to ZAFramework", map[string]any{})
}

func (c *Controller) Settings(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/wrapper/settings.html", "Settings", map[string]any{})
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Hapus data spesifik dari session (misalnya token)
	// Jika library session Anda mendukung Delete, gunakan itu.
	// Jika tidak, mengosongkan nilainya sudah cukup untuk memutus autentikasi.

	err := ehttp.DestroyAuthorize(r)
	if err != nil {
		fmt.Println("gagal hapus sesi di redis", err)
	}

	_, err = session.Set(w, r, "token", "")
	if err != nil {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Gagal membersihkan sesi"})
		return
	}

	// 2. Hancurkan seluruh session (Hapus Cookie dari browser)
	// Biasanya library session memiliki fungsi Destroy atau Clear
	err = session.Destroy(w, r)
	if err != nil {
		// Tetap lanjut meskipun destroy gagal di server,
		// yang penting di sisi client token sudah kosong.
		fmt.Println("Gagal destroy session:", err)
	}

	// 3. Respon sesuai kebutuhan
	// Jika diakses via AJAX (SPA), kirim JSON
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" || strings.Contains(r.Header.Get("Accept"), "application/json") {
		c.Response.JSON(w, r, map[string]any{
			"status": true,
			"msg":    "Anda telah keluar dari sistem.",
		})
		return
	}

	// Jika diakses via link biasa, langsung redirect ke login
	http.Redirect(w, r, utils.GetEnv("base_url"), http.StatusSeeOther)
}
