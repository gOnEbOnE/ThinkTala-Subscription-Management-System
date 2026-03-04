package notification

import "time"

// Notification adalah model untuk broadcast notification yang ditampilkan ke user/client.
type Notification struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`        // info | warning | error | success
	TargetRole string    `json:"target_role"` // all | client | ops | compliance
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  *string   `json:"created_by,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateNotificationRequest adalah payload untuk membuat notification baru.
type CreateNotificationRequest struct {
	Title      string `json:"title"   binding:"required"`
	Message    string `json:"message" binding:"required"`
	Type       string `json:"type"`
	TargetRole string `json:"target_role"`
	IsActive   *bool  `json:"is_active"`
	CreatedBy  string `json:"created_by"`
}

// UpdateNotificationRequest adalah payload untuk update notification.
type UpdateNotificationRequest struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	TargetRole string `json:"target_role"`
	IsActive   *bool  `json:"is_active"`
	UpdatedBy  string `json:"updated_by"`
}
