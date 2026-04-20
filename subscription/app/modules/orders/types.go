package orders

import "time"

// CreateOrderDTO - request body untuk POST /api/orders (PBI-45)
type CreateOrderDTO struct {
	PackageID      string `json:"package_id"`
	DurationMonths int    `json:"duration_months"` // opsional, default 1 bulan
	PaymentMethod  string `json:"payment_method"`
}

// VerifyOrderDTO - request body untuk PATCH /api/admin/orders/{id}/verify
type VerifyOrderDTO struct {
	Action       string `json:"action"`
	RejectReason string `json:"reject_reason,omitempty"`
}

// Order - mapping dari tabel subscription.orders
type Order struct {
	ID                     string     `json:"order_id"`
	InvoiceNumber          string     `json:"invoice_number"`
	UserID                 string     `json:"user_id,omitempty"`
	PackageID              string     `json:"package_id"`
	DurationMonths         int        `json:"duration_months"`
	PaymentMethod          string     `json:"payment_method"`
	TotalPrice             float64    `json:"total_price"`
	Status                 string     `json:"status"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// ClientOrderListItem - item riwayat pesanan milik client
type ClientOrderListItem struct {
	OrderID                string     `json:"order_id"`
	InvoiceNumber          string     `json:"invoice_number"`
	PackageName            string     `json:"package_name"`
	TotalPrice             float64    `json:"total_price"`
	PaymentMethod          string     `json:"payment_method"`
	Status                 string     `json:"status"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// ClientOrderDetail - detail pesanan untuk endpoint GET /api/orders/{id}
type ClientOrderDetail struct {
	OrderID                string     `json:"order_id"`
	InvoiceNumber          string     `json:"invoice_number"`
	PackageName            string     `json:"package_name"`
	TotalPrice             float64    `json:"total_price"`
	PaymentMethod          string     `json:"payment_method"`
	Status                 string     `json:"status"`
	VerificationNote       string     `json:"verification_note,omitempty"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	PaymentProofURL        string     `json:"payment_proof_url,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// AdminOrderListItem - item daftar pesanan untuk operasional
type AdminOrderListItem struct {
	OrderID                string     `json:"order_id"`
	InvoiceNumber          string     `json:"invoice_number"`
	ClientName             string     `json:"client_name"`
	PackageName            string     `json:"package_name"`
	TotalPrice             float64    `json:"total_price"`
	PaymentMethod          string     `json:"payment_method"`
	Status                 string     `json:"status"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// AdminOrderDetail - detail pesanan untuk verifikasi operasional
type AdminOrderDetail struct {
	OrderID                string     `json:"order_id"`
	InvoiceNumber          string     `json:"invoice_number"`
	UserID                 string     `json:"user_id"`
	ClientName             string     `json:"client_name"`
	ClientEmail            string     `json:"client_email"`
	PackageID              string     `json:"package_id"`
	PackageName            string     `json:"package_name"`
	DurationMonths         int        `json:"duration_months"`
	TotalPrice             float64    `json:"total_price"`
	PaymentMethod          string     `json:"payment_method"`
	Status                 string     `json:"status"`
	VerificationNote       string     `json:"verification_note,omitempty"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	PaymentProofURL        string     `json:"payment_proof_url,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

// VerifyResult - response untuk endpoint verifikasi pembayaran
type VerifyResult struct {
	Message          string `json:"message"`
	OrderID          string `json:"order_id"`
	NewStatus        string `json:"new_status"`
	VerificationNote string `json:"verification_note,omitempty"`
}

// ActivationResult - response untuk endpoint aktivasi subscription
type ActivationResult struct {
	Message        string    `json:"message"`
	SubscriptionID string    `json:"subscription_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
}

// SubscriptionStatus - ringkasan subscription aktif untuk dashboard client
type SubscriptionStatus struct {
	SubscriptionID string    `json:"subscription_id"`
	PackageID      string    `json:"package_id,omitempty"`
	PackageName    string    `json:"package_name"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
}

// UploadPaymentProofResult - response upload bukti transfer
type UploadPaymentProofResult struct {
	OrderID                string     `json:"order_id"`
	HasPaymentProof        bool       `json:"has_payment_proof"`
	PaymentProofUploadedAt *time.Time `json:"payment_proof_uploaded_at,omitempty"`
	Message                string     `json:"message"`
}

// PaymentProofFile - representasi file bukti transfer
type PaymentProofFile struct {
	FileName    string
	ContentType string
	Data        []byte
	UploadedAt  *time.Time
}

// OrderRecord - representasi internal order lengkap
type OrderRecord struct {
	OrderID                 string
	InvoiceNumber           string
	UserID                  string
	PackageID               string
	PackageName             string
	ClientName              string
	ClientEmail             string
	DurationMonths          int
	TotalPrice              float64
	PaymentMethod           string
	Status                  string
	VerificationNote        string
	HasPaymentProof         bool
	PaymentProofUploadedAt  *time.Time
	PaymentProofFileName    string
	PaymentProofContentType string
	CreatedAt               time.Time
}

// PackageInfo - subset data paket untuk validasi + pricing lookup
type PackageInfo struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Price  float64 `json:"price"` // harga dasar / fallback
}
