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
