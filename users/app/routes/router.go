package routes

import (
	"net/http"
	"time"

	"github.com/master-abror/zaframework/app/modules/admin"
	"github.com/master-abror/zaframework/app/modules/kyc"
	"github.com/master-abror/zaframework/app/modules/login"
	"github.com/master-abror/zaframework/app/modules/register"
	"github.com/master-abror/zaframework/core"
	middleware "github.com/master-abror/zaframework/core/http"
)

// Init mendaftarkan semua route aplikasi
func Init(app *core.App,
	loginController *login.Controller,
	registerController *register.Controller,
	kycController *kyc.Controller,
	kycAdminController *kyc.AdminController,
	adminController *admin.Controller,
) {

	// ==============================
	// 1. Static Assets
	// ==============================
	fs := http.FileServer(http.Dir("./public/assets"))
	app.Router.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	// ==============================
	// 2. Root redirect to login
	// ==============================
	app.Router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/account/login", http.StatusFound)
	})

	// ==============================
	// 3. System Routes (Health Check)
	// ==============================
	app.Router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		app.Response.JSON(w, r, map[string]any{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// ==============================
	// 4. Auth Routes (Pages)
	// ==============================
	app.Router.HandleFunc("GET /account/login", loginController.Login)
	app.Router.HandleFunc("POST /account/login/auth", loginController.Auth)
	app.Router.HandleFunc("GET /account/register", registerController.Register)
	app.Router.HandleFunc("GET /account/verify-otp", registerController.VerifyOTPPage)

	// ==============================
	// 5. API Routes (Protected)
	// ==============================
	app.Router.Handle("POST /api/auth", middleware.ApiMiddleware(loginController.ApiAuth))

	// ==============================
	// 6. API Routes (Public - No Auth Required)
	// ==============================
	app.Router.HandleFunc("POST /api/auth/register", registerController.Submit)
	app.Router.HandleFunc("POST /api/auth/verify-otp", registerController.VerifyOTP)
	app.Router.HandleFunc("POST /api/auth/resend-otp", registerController.ResendOTP)
	app.Router.HandleFunc("POST /api/auth/logout", loginController.Logout)
	app.Router.HandleFunc("POST /api/auth/assume-role", loginController.AssumeRole)
	app.Router.HandleFunc("GET /api/auth/roles", loginController.GetRoles)

	// ==============================
	// 7. KYC Routes (Protected via Gateway Auth)
	// ==============================
	app.Router.HandleFunc("POST /api/kyc/submit", kycController.Submit)
	app.Router.HandleFunc("GET /api/kyc/status", kycController.Status)

	// ==============================
	// 8. Admin KYC Routes (Protected via Gateway Auth)
	// ==============================
	app.Router.Handle("GET /api/admin/kyc", http.HandlerFunc(kycAdminController.ServeHTTP))
	app.Router.Handle("GET /api/admin/kyc/{id}", http.HandlerFunc(kycAdminController.ServeHTTP))
	app.Router.Handle("POST /api/admin/kyc/{id}/approve", http.HandlerFunc(kycAdminController.ServeHTTP))
	app.Router.Handle("POST /api/admin/kyc/{id}/reject", http.HandlerFunc(kycAdminController.ServeHTTP))

	// ==============================
	// 9. Admin User Management Routes (Protected via Gateway Auth - SUPERADMIN only)
	// ==============================
	app.Router.HandleFunc("POST /api/admin/users", adminController.CreateUser)

}
