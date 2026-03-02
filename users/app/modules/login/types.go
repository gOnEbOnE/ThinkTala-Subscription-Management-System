package login

import "time"

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Photo     *string   `json:"photo"`
	GroupID   string    `json:"group_id"`
	LevelID   string    `json:"level_id"`
	RoleID    string    `json:"role_id"`
	Status    string    `json:"status"`
	Level     string    `json:"level"`
	LevelCode string    `json:"level_code"`
	Role      string    `json:"role"`
	RoleCode  string    `json:"role_code"`
	Group     string    `json:"group"`
	GroupCode string    `json:"group_code"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}

// Struct untuk payload input dari Controller -> Service
type LoginPayload struct {
	LoginID   string `json:"login_id"`
	Password  string `json:"password"`
	IP        string `json:"ip"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Struct output Service -> Controller
type LoginResult struct {
	Token    string         `json:"token"`
	UserData map[string]any `json:"user_data"`
}
