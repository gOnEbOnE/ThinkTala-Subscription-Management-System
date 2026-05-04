package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/master-abror/zaframework/core/database"
)

type Handler struct {
	DB *database.DBWrapper
}

func NewHandler(db *database.DBWrapper) *Handler {
	return &Handler{DB: db}
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

type orderStats struct {
	TotalPendingPayment int `json:"total_pending_payment"`
	TotalPaid           int `json:"total_paid"`
	TotalCancelled      int `json:"total_cancelled"`
	TotalToday          int `json:"total_today"`
}

type activeSubStats struct {
	TotalActive  int `json:"total_active"`
	TotalExpired int `json:"total_expired"`
}

type recentOrder struct {
	ID            string  `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	UserID        string  `json:"user_id"`
	Status        string  `json:"status"`
	TotalPrice    float64 `json:"total_price"`
	CreatedAt     string  `json:"created_at"`
}

// GetOperationalDashboard handles GET /api/superadmin/dashboard/operational
func (h *Handler) GetOperationalDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}

	ctx := r.Context()
	oStats := orderStats{}

	row := h.DB.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'PENDING_PAYMENT') AS pending_payment,
			COUNT(*) FILTER (WHERE status = 'PAID')            AS paid,
			COUNT(*) FILTER (WHERE status = 'CANCELLED')       AS cancelled,
			COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE) AS today
		FROM subscription.orders
	`)
	if err := row.Scan(&oStats.TotalPendingPayment, &oStats.TotalPaid, &oStats.TotalCancelled, &oStats.TotalToday); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil statistik order"})
		return
	}

	subStats := activeSubStats{}
	row = h.DB.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'ACTIVE')  AS active,
			COUNT(*) FILTER (WHERE status = 'EXPIRED') AS expired
		FROM subscription.subscriptions
	`)
	if err := row.Scan(&subStats.TotalActive, &subStats.TotalExpired); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil statistik subscription"})
		return
	}

	rows, err := h.DB.Pool.Query(ctx, `
		SELECT id, invoice_number, user_id::text, status, total_price, created_at::text
		FROM subscription.orders
		WHERE status = 'PENDING_PAYMENT'
		ORDER BY created_at ASC
		LIMIT 10
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil order pending"})
		return
	}
	defer rows.Close()

	recent := []recentOrder{}
	for rows.Next() {
		var o recentOrder
		if err := rows.Scan(&o.ID, &o.InvoiceNumber, &o.UserID, &o.Status, &o.TotalPrice, &o.CreatedAt); err != nil {
			continue
		}
		recent = append(recent, o)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"order_stats":        oStats,
		"subscription_stats": subStats,
		"recent_pending":     recent,
	})
}

// GetInternalOpsSummary handles GET /internal/dashboard/ops-summary (for B8 overview)
func (h *Handler) GetInternalOpsSummary(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil || h.DB.Pool == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}

	ctx := r.Context()
	var pendingPayment, paid, activeSubscriptions int

	row := h.DB.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'PENDING_PAYMENT') AS pending_payment,
			COUNT(*) FILTER (WHERE status = 'PAID')            AS paid
		FROM subscription.orders
	`)
	if err := row.Scan(&pendingPayment, &paid); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
		return
	}

	if err := h.DB.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM subscription.subscriptions WHERE status = 'ACTIVE'
	`).Scan(&activeSubscriptions); err != nil {
		activeSubscriptions = 0
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pending_payment":      pendingPayment,
		"paid":                 paid,
		"active_subscriptions": activeSubscriptions,
	})
}
