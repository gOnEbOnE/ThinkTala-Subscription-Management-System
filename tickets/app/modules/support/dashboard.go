package support

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type ticketStats struct {
	TotalOpen             int     `json:"total_open"`
	TotalResolved         int     `json:"total_resolved"`
	TotalToday            int     `json:"total_today"`
	TotalAll              int     `json:"total_all"`
	OverallCompletionRate float64 `json:"overall_completion_rate"`
}

type adminSupportStat struct {
	AdminID         string  `json:"admin_id"`
	TicketsHandled  int     `json:"tickets_handled"`
	CompletionRate  float64 `json:"completion_rate"`
	TotalReplies    int     `json:"total_replies"`
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
				COUNT(*) FILTER (WHERE status != 'DONE')  AS open,
				COUNT(*) FILTER (WHERE status = 'DONE')   AS resolved,
				COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE) AS today,
				COUNT(*)                                   AS total,
				COALESCE(ROUND(
					COUNT(*) FILTER (WHERE status = 'DONE')::numeric / NULLIF(COUNT(*), 0) * 100, 1
				), 0) AS overall_completion_rate
			FROM support_tickets
		`)
		if err := row.Scan(&stats.TotalOpen, &stats.TotalResolved, &stats.TotalToday, &stats.TotalAll, &stats.OverallCompletionRate); err != nil {
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

		adminRows, adminErr := db.QueryContext(ctx, `
			SELECT
				str.admin_id,
				COUNT(DISTINCT str.ticket_id) AS tickets_handled,
				COALESCE(ROUND(
					COUNT(DISTINCT CASE WHEN st.status != 'ON PROCESS' THEN str.ticket_id END)::numeric
					/ NULLIF(COUNT(DISTINCT str.ticket_id), 0) * 100, 1
				), 0) AS completion_rate,
				COUNT(*) AS total_replies
			FROM support_ticket_replies str
			JOIN support_tickets st ON st.id = str.ticket_id
			GROUP BY str.admin_id
			ORDER BY tickets_handled DESC
		`)
		adminStats := []adminSupportStat{}
		if adminErr == nil {
			defer adminRows.Close()
			for adminRows.Next() {
				var s adminSupportStat
				if scanErr := adminRows.Scan(&s.AdminID, &s.TicketsHandled, &s.CompletionRate, &s.TotalReplies); scanErr == nil {
					adminStats = append(adminStats, s)
				}
			}
		}

		writeDashJSON(w, http.StatusOK, map[string]any{
			"stats":        stats,
			"recent_open":  recent,
			"admin_stats":  adminStats,
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
