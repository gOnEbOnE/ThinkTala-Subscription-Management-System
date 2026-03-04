package routes

import (
	notifmod "notification/app/modules/notification_broadcast"
	tplmod "notification/app/modules/template_notification"

	"github.com/gin-gonic/gin"
)

// Register mendaftarkan semua routes ke Gin engine.
// Dipanggil sekali dari main.go saat service startup.
func Register(r *gin.Engine) {
	// ── Static Pages ─────────────────────────────────────────────
	r.GET("/compliance/dashboard", func(c *gin.Context) {
		c.File("./public/views/dashboard.html")
	})

	// ── Broadcast Notifications ───────────────────────────────────
	// Digunakan oleh ops dashboard untuk membuat announcement/banner.
	notifCtrl := notifmod.NewController()
	api := r.Group("/api/notifications")
	{
		api.GET("", notifCtrl.List)
		api.GET("/public", notifCtrl.ListPublic) // harus sebelum /:id agar tidak collision
		api.GET("/:id", notifCtrl.Get)
		api.POST("", notifCtrl.Create)
		api.PUT("/:id", notifCtrl.Update)
		api.DELETE("/:id", notifCtrl.Delete)
	}

	// ── Notification Templates ────────────────────────────────────
	// Digunakan untuk mendefinisikan template pesan per event_type & channel.
	tplCtrl := tplmod.NewController()
	tpl := r.Group("/api/help/notification-templates")
	{
		tpl.GET("", tplCtrl.List)
		tpl.GET("/:id", tplCtrl.Get)
		tpl.POST("", tplCtrl.Create)
		tpl.PUT("/:id", tplCtrl.Update)
		tpl.DELETE("/:id", tplCtrl.Delete)
	}
}
