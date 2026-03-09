package routes

import (
	notifmod "notification/app/modules/notification_broadcast"
	tplmod "notification/app/modules/template_notification"

	"github.com/gin-gonic/gin"
)

// Register mendaftarkan semua routes ke Gin engine dan mengembalikan Service untuk dipakai worker.
// Dipanggil sekali dari main.go saat service startup.
func Register(r *gin.Engine) *tplmod.Service {
	tplSvc := tplmod.NewService()
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
	tplCtrl := tplmod.NewController(tplSvc)
	tpl := r.Group("/api/help/notification-templates")
	{
		tpl.GET("", tplCtrl.List)
		tpl.GET("/:id", tplCtrl.Get)
		tpl.POST("", tplCtrl.Create)
		tpl.PUT("/:id", tplCtrl.Update)
		tpl.DELETE("/:id", tplCtrl.Delete)
	}
	// POST /api/notifications/send — kirim notifikasi berdasarkan event_type + template
	api.POST("/send", tplCtrl.Send)

	// GET /api/help/notification-templates/event-types — daftar event_type yang sudah terdaftar
	tpl.GET("/event-types", tplCtrl.EventTypes)

	// GET /api/notifications/logs — monitoring log pengiriman (dengan filter status opsional)
	api.GET("/logs", tplCtrl.Logs)

	return tplSvc
}
