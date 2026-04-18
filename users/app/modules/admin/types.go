package admin

import "time"

// AllowedRoles defines the valid roles that can be assigned by Superadmin
var AllowedRoles = map[string]bool{
	"OPERASIONAL": true,
	"COMPLIANCE":  true,
	"MANAJEMEN":   true,
	"ADMIN_CS":    true,
}

// CreateUserInput — payload dari frontend untuk membuat user internal
type CreateUserInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// CreateUserResult — output setelah user berhasil dibuat (tanpa password)
type CreateUserResult struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UserListItem — data user untuk daftar akun internal (PBI-52)
type UserListItem struct {
	UserID      string     `json:"user_id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// GetUsersParams — parameter filter dan paginasi (PBI-52)
type GetUsersParams struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Role    string `json:"role"`
	Status  string `json:"status"`
	Search  string `json:"search"`
}

// GetUsersResponse — format output yang di-pass ke ApiJSON data field (PBI-52)
// Matches spec: {"data": [...], "total": N, "page": N, "per_page": N}
type GetUsersResponse struct {
	Data    []UserListItem `json:"data"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

// UserDetail — data lengkap satu user untuk detail page (PBI-53)
type UserDetail struct {
	UserID      string     `json:"user_id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
