package orders

import "time"

// ==========================================
// 1. DATA TRANSFER OBJECTS (DTO)
// ==========================================

// CreateOrderDTO - request body untuk POST /api/orders (PBI-45)
type CreateOrderDTO struct {
	PackageID     string `json:"package_id"`
	PaymentMethod string `json:"payment_method"`
}

// ==========================================
// 2. MODEL STRUCTS
// ==========================================

// Order - mapping row dari tabel subscription.orders
type Order struct {
	ID            string    `json:"order_id"`
	InvoiceNumber string    `json:"invoice_number"`
	UserID        string    `json:"user_id,omitempty"`
	PackageID     string    `json:"package_id"`
	PaymentMethod string    `json:"payment_method"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// PackageInfo - subset data paket yang dibutuhkan saat buat order
type PackageInfo struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}
