package compliance

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/master-abror/zaframework/core/database"
)

type DashboardHandler struct {
	DB *database.DBWrapper
}

func NewDashboardHandler(db *database.DBWrapper) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

type kycStats struct {
	TotalPending  int `json:"total_pending"`
	TotalApproved int `json:"total_approved"`
	TotalRejected int `json:"total_rejected"`
	TotalToday    int `json:"total_today"`
}

type recentItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	NIK       string    `json:"nik"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetDashboard handles GET /api/superadmin/dashboard/compliance
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}

	ctx := r.Context()
	stats := kycStats{}

	row := h.DB.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending')  AS pending,
			COUNT(*) FILTER (WHERE status = 'approved') AS approved,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected,
			COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE) AS today
		FROM kyc_submissions
	`)
	if err := row.Scan(&stats.TotalPending, &stats.TotalApproved, &stats.TotalRejected, &stats.TotalToday); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil statistik KYC"})
		return
	}

	rows, err := h.DB.Pool.Query(ctx, `
		SELECT ks.id, ks.user_id, ks.full_name, ks.nik,
		       COALESCE(u.email, '') AS email, ks.status, ks.created_at
		FROM kyc_submissions ks
		LEFT JOIN users u ON ks.user_id::text = u.id::text
		WHERE ks.status = 'pending'
		ORDER BY ks.created_at ASC
		LIMIT 10
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data KYC pending"})
		return
	}
	defer rows.Close()

	recent := []recentItem{}
	for rows.Next() {
		var item recentItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.FullName, &item.NIK, &item.Email, &item.Status, &item.CreatedAt); err != nil {
			continue
		}
		recent = append(recent, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"stats":          stats,
		"recent_pending": recent,
	})
}

// GetInternalSummary handles GET /internal/dashboard/kyc-summary (for B8 overview aggregation)
func (h *DashboardHandler) GetInternalSummary(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}

	ctx := r.Context()
	var pending, approved, rejected int

	row := h.DB.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending')  AS pending,
			COUNT(*) FILTER (WHERE status = 'approved') AS approved,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected
		FROM kyc_submissions
	`)
	if err := row.Scan(&pending, &approved, &rejected); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pending":  pending,
		"approved": approved,
		"rejected": rejected,
	})
}
