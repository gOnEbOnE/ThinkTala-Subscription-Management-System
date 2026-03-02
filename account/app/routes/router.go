package routes

import (
	"net/http"
	"time"

	"github.com/master-abror/zaframework/app/modules/account"
	"github.com/master-abror/zaframework/app/modules/dashboard"
	"github.com/master-abror/zaframework/core"
	middleware "github.com/master-abror/zaframework/core/http"
)

// Init mendaftarkan semua route aplikasi
func Init(app *core.App,
	accountController *account.Controller,
	dashboardController *dashboard.Controller,
) {

	// ==============================
	// 1. Static Assets
	// ==============================
	// Mengakses folder ./assets via browser
	fs := http.FileServer(http.Dir("./public/assets"))
	app.Router.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	// ==============================
	// 2. System Routes (Health Check)
	// ==============================
	app.Router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		app.Response.JSON(w, r, map[string]any{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	app.Router.Handle("GET /account", middleware.UserAuthorize(accountController.Index))
	app.Router.Handle("GET /account/dashboard", middleware.UserAuthorize(dashboardController.Index))
	app.Router.HandleFunc("GET /account/logout", accountController.Logout)

}
