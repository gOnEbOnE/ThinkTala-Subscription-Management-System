package routes

import (
	"net/http"
	"time"

	"github.com/master-abror/zaframework/app/modules/orders"
	"github.com/master-abror/zaframework/app/modules/packages"
	"github.com/master-abror/zaframework/core"
)

// Init mendaftarkan semua route subscription service
func Init(app *core.App, packagesController *packages.Controller, ordersController *orders.Controller) {

	// Static assets
	fs := http.FileServer(http.Dir("./public/assets"))
	app.Router.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	// Health Check
	app.Router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		app.Response.JSON(w, r, map[string]any{
			"status":  "healthy",
			"service": "subscription",
			"time":    time.Now(),
		})
	})

	// ===================================================
	// EPIC04: Package Management
	// ===================================================

	// PBI-32: Create Package (Admin)
	app.Router.HandleFunc("POST /api/admin/packages", packagesController.CreatePackageHandler)

	// PBI-33 + PBI-36: Get Packages for Admin (with filter)
	app.Router.HandleFunc("GET /api/admin/packages", packagesController.GetPackagesAdminHandler)

	// PBI-34: Update Package (Admin)
	app.Router.HandleFunc("PUT /api/admin/packages/{id}", packagesController.UpdatePackageHandler)

	// PBI-35: Delete Package / Soft Delete (Admin)
	app.Router.HandleFunc("DELETE /api/admin/packages/{id}", packagesController.DeletePackageHandler)

	// Toggle Package Status (Activate / Deactivate)
	app.Router.HandleFunc("PATCH /api/admin/packages/{id}/status", packagesController.TogglePackageStatusHandler)

	// PBI-37: Get Active Catalog for Client
	app.Router.HandleFunc("GET /api/subscription/catalog", packagesController.GetCatalogHandler)

	// Legacy alias sesuai routes.json gateway
	app.Router.HandleFunc("GET /api/subscriptions", packagesController.GetPackagesAdminHandler)

	// ===================================================
	// EPIC04 PBI-45: Create Order (Client)
	// ===================================================
	app.Router.HandleFunc("POST /api/orders", ordersController.CreateOrderHandler)
	app.Router.HandleFunc("GET /api/orders", ordersController.ListOrdersClientHandler)
	app.Router.HandleFunc("GET /api/orders/{id}", ordersController.GetOrderDetailClientHandler)
	app.Router.HandleFunc("POST /api/orders/{id}/payment-proof", ordersController.UploadPaymentProofClientHandler)
	app.Router.HandleFunc("GET /api/orders/{id}/payment-proof", ordersController.GetPaymentProofClientHandler)

	// ===================================================
	// EPIC04 PBI-49: Verify Order (Operasional)
	// ===================================================
	app.Router.HandleFunc("GET /api/admin/orders", ordersController.ListOrdersAdminHandler)
	app.Router.HandleFunc("GET /api/admin/orders/{id}", ordersController.GetOrderDetailAdminHandler)
	app.Router.HandleFunc("GET /api/admin/orders/{id}/payment-proof", ordersController.GetPaymentProofAdminHandler)
	app.Router.HandleFunc("PATCH /api/admin/orders/{id}/verify", ordersController.VerifyOrderHandler)

	// ===================================================
	// EPIC04 PBI-50: Activate Subscription after Verify
	// ===================================================
	app.Router.HandleFunc("PATCH /api/admin/orders/{id}/activate", ordersController.ActivateOrderHandler)
	app.Router.HandleFunc("GET /api/subscriptions/me", ordersController.GetMySubscriptionHandler)
}
