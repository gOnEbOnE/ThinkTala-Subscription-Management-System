package notification

import "time"

// Notification adalah model untuk broadcast notification yang ditampilkan ke user/client.
type Notification struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Message     string     `json:"message,omitempty"` // backward-compat alias of description
	Type        string     `json:"type"`
	CTAURL      *string    `json:"cta_url,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	ExpiryDate  *time.Time `json:"expiry_date,omitempty"`
	Status      string     `json:"status"` // active | inactive | expired
	TargetRole  string     `json:"target_role,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *string    `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *string    `json:"updated_by,omitempty"`
}

// CreateNotificationRequest adalah payload untuk membuat notification baru.
type CreateNotificationRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Message     string `json:"message,omitempty"` // backward-compat alias of description
	Type        string `json:"type"`
	CTAURL      string `json:"cta_url"`
	ImageURL    string `json:"image_url"`
	ExpiryDate  string `json:"expiry_date"`
	TargetRole  string `json:"target_role,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
}

// UpdateNotificationRequest adalah payload untuk update notification.
type UpdateNotificationRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"` // backward-compat alias of description
	Type        string `json:"type,omitempty"`
	CTAURL      string `json:"cta_url,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	ExpiryDate  string `json:"expiry_date,omitempty"`
	TargetRole  string `json:"target_role,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
	UpdatedBy   string `json:"updated_by,omitempty"`
}
