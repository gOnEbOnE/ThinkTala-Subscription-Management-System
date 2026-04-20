package routes

import (
	"encoding/json"
	"net/http"

	"tickets/app/modules/support"
	"tickets/core"
)

func Init(app *core.App) {
	app.Router.HandleFunc("/api/admin/support/tickets", support.HandleAdminSupportTickets(app.DB))
	app.Router.HandleFunc("/api/admin/support/tickets/", support.HandleAdminSupportTickets(app.DB))
	app.Router.HandleFunc("/api/admin/support/tickets/attachment", support.HandleAdminSupportTicketAttachment(app.DB))
	app.Router.HandleFunc("/api/admin/support/tickets/attachment/", support.HandleAdminSupportTicketAttachment(app.DB))
	app.Router.HandleFunc("/api/user/support/tickets", support.HandleCreateUserSupportTicket(app.DB))
	app.Router.HandleFunc("/api/user/support/tickets/", support.HandleCreateUserSupportTicket(app.DB))

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": "Tickets Service",
			"status":  "Healthy",
			"path":    r.URL.Path,
		})
	})
}
