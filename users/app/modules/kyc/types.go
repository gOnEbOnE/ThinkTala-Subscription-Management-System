package kyc

import "time"

// KYCSubmission — model database kyc_submissions
type KYCSubmission struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	FullName       string     `json:"full_name"`
	NIK            string     `json:"nik"`
	Address        string     `json:"address"`
	Birthdate      string     `json:"birthdate"`
	Phone          string     `json:"phone"`
	KTPImage       string     `json:"ktp_image"`
	Status         string     `json:"status"`
	RejectReason   *string    `json:"reject_reason,omitempty"`
	RejectedFields []string   `json:"rejected_fields,omitempty"`
	ReviewedBy     *string    `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// KYCSubmitPayload — payload dari controller ke service (via dispatcher)
type KYCSubmitPayload struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	NIK       string `json:"nik"`
	Address   string `json:"address"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone"`
	KTPImage  string `json:"ktp_image"` // filename hasil upload
}

// KYCSubmitResult — output dari service ke controller
type KYCSubmitResult struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// KYCStatusResult — output untuk cek status KYC
type KYCStatusResult struct {
	ID             string     `json:"id"`
	FullName       string     `json:"full_name"`
	NIK            string     `json:"nik"`
	Address        string     `json:"address"`
	Birthdate      string     `json:"birthdate"`
	Phone          string     `json:"phone"`
	KTPImage       string     `json:"ktp_image"`
	Status         string     `json:"status"`
	RejectReason   *string    `json:"reject_reason,omitempty"`
	RejectedFields []string   `json:"rejected_fields,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ========== ADMIN KYC TYPES ==========

// KYCListItem — item untuk daftar KYC di admin panel
type KYCListItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	NIK       string    `json:"nik"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// KYCDetailResult — detail lengkap KYC untuk admin review
type KYCDetailResult struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	FullName       string     `json:"full_name"`
	NIK            string     `json:"nik"`
	Address        string     `json:"address"`
	Birthdate      string     `json:"birthdate"`
	Phone          string     `json:"phone"`
	KTPImage       string     `json:"ktp_image"`
	Email          string     `json:"email"`
	Status         string     `json:"status"`
	RejectReason   *string    `json:"reject_reason,omitempty"`
	RejectedFields []string   `json:"rejected_fields,omitempty"`
	ReviewedBy     *string    `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// KYCReviewPayload — payload approve/reject dari controller ke service
type KYCReviewPayload struct {
	KYCID          string   `json:"kyc_id"`
	ReviewerID     string   `json:"reviewer_id"`
	Action         string   `json:"action"` // "approve" atau "reject"
	RejectReason   string   `json:"reject_reason,omitempty"`
	RejectedFields []string `json:"rejected_fields,omitempty"`
}
