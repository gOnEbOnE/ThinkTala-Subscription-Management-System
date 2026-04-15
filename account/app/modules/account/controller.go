package account

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	// Import package http core kamu
	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/session"
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

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("current_user").(*token.CurrentUser)

	fmt.Println("user", user)

	data := map[string]any{
		"AppName": "ZaFramework",
		"Year":    time.Now().Year(),
	}

	c.Response.View(w, r, "public/views/account/page.html", "Account", data)
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {

	err := ehttp.DestroyAuthorize(r)
	if err != nil {
		fmt.Println("gagal hapus sesi di redis", err)
	}

	_, err = session.Set(w, r, "token", "")
	if err != nil {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Gagal membersihkan sesi"})
		return
	}

	err = session.Destroy(w, r)
	if err != nil {
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

	http.Redirect(w, r, utils.GetEnv("base_url"), http.StatusSeeOther)
}
