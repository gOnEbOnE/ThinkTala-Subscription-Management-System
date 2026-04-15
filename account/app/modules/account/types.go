package account

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
	Role      string    `json:"role"`
	Group     string    `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}

type CurrentUser struct {
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
	Role      string    `json:"role"`
	Group     string    `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}
