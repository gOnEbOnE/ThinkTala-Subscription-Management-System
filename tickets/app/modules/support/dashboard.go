package support

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type ticketStats struct {
	TotalOpen     int `json:"total_open"`
	TotalResolved int `json:"total_resolved"`
	TotalToday    int `json:"total_today"`
	TotalAll      int `json:"total_all"`
}

type recentTicket struct {
	ID             string    `json:"id"`
	ReporterName   string    `json:"reporter_name"`
	ReporterEmail  string    `json:"reporter_email"`
	Title          string    `json:"title"`
	Category       string    `json:"category"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func writeDashJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// HandleSupportDashboard handles GET /api/superadmin/dashboard/support
func HandleSupportDashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeDashJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

		if db == nil {
			writeDashJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
			return
		}

		ctx := r.Context()
		stats := ticketStats{}

		row := db.QueryRowContext(ctx, `
			SELECT
				COUNT(*) FILTER (WHERE status != 'RESOLVED')  AS open,
				COUNT(*) FILTER (WHERE status = 'RESOLVED')   AS resolved,
				COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE) AS today,
				COUNT(*)                                       AS total
			FROM support_tickets
		`)
		if err := row.Scan(&stats.TotalOpen, &stats.TotalResolved, &stats.TotalToday, &stats.TotalAll); err != nil {
			writeDashJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil statistik tiket"})
			return
		}

		rows, err := db.QueryContext(ctx, `
			SELECT id, reporter_name, reporter_email, title, category, status, created_at
			FROM support_tickets
			WHERE status != 'RESOLVED'
			ORDER BY created_at ASC
			LIMIT 10
		`)
		if err != nil {
			writeDashJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil tiket terbuka"})
			return
		}
		defer rows.Close()

		recent := []recentTicket{}
		for rows.Next() {
			var t recentTicket
			if err := rows.Scan(&t.ID, &t.ReporterName, &t.ReporterEmail, &t.Title, &t.Category, &t.Status, &t.CreatedAt); err != nil {
				continue
			}
			recent = append(recent, t)
		}

		writeDashJSON(w, http.StatusOK, map[string]any{
			"stats":        stats,
			"recent_open":  recent,
		})
	}
}

// HandleSupportInternalSummary handles GET /internal/dashboard/support-summary (for B8 overview)
func HandleSupportInternalSummary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeDashJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
			return
		}

		ctx := r.Context()
		var open, resolved int

		row := db.QueryRowContext(ctx, `
			SELECT
				COUNT(*) FILTER (WHERE status != 'RESOLVED') AS open,
				COUNT(*) FILTER (WHERE status = 'RESOLVED')  AS resolved
			FROM support_tickets
		`)
		if err := row.Scan(&open, &resolved); err != nil {
			writeDashJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
			return
		}

		writeDashJSON(w, http.StatusOK, map[string]any{
			"open":     open,
			"resolved": resolved,
		})
	}
}
