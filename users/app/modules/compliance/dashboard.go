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

type adminKYCStat struct {
	AdminID        string  `json:"admin_id"`
	AdminName      string  `json:"admin_name"`
	TotalProcessed int     `json:"total_processed"`
	Approved       int     `json:"approved"`
	Rejected       int     `json:"rejected"`
	SLARate        float64 `json:"sla_rate"`
	AvgReviewHours float64 `json:"avg_review_hours"`
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

	adminRows, err := h.DB.Pool.Query(ctx, `
		SELECT
			u.id AS admin_id,
			COALESCE(NULLIF(u.email, ''), u.id::text) AS admin_name,
			COUNT(ks.id) AS total_processed,
			COUNT(*) FILTER (WHERE ks.status = 'approved') AS approved,
			COUNT(*) FILTER (WHERE ks.status = 'rejected') AS rejected,
			COALESCE(ROUND(
				COUNT(*) FILTER (WHERE ks.reviewed_at IS NOT NULL AND ks.reviewed_at - ks.created_at < INTERVAL '24 hours')::numeric
				/ NULLIF(COUNT(*), 0) * 100, 1
			), 0) AS sla_rate,
			COALESCE(ROUND(
				EXTRACT(EPOCH FROM AVG(ks.reviewed_at - ks.created_at)) / 3600, 1
			), 0) AS avg_review_hours
		FROM kyc_submissions ks
		JOIN users u ON u.id::text = ks.reviewed_by::text
		WHERE ks.reviewed_by IS NOT NULL
		GROUP BY u.id, u.email
		ORDER BY total_processed DESC
	`)
	adminStats := []adminKYCStat{}
	if err == nil {
		defer adminRows.Close()
		for adminRows.Next() {
			var s adminKYCStat
			if scanErr := adminRows.Scan(
				&s.AdminID, &s.AdminName,
				&s.TotalProcessed, &s.Approved, &s.Rejected,
				&s.SLARate, &s.AvgReviewHours,
			); scanErr == nil {
				adminStats = append(adminStats, s)
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"stats":          stats,
		"recent_pending": recent,
		"admin_stats":    adminStats,
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
