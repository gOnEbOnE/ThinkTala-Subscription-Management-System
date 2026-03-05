package register

import "time"

// RegisterInput — payload dari frontend
type RegisterInput struct {
	FullName  string `json:"full_name" validate:"required,min=3,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Phone     string `json:"no_telp" validate:"required"`
	Birthdate string `json:"tanggal_lahir" validate:"required"` // Format: YYYY-MM-DD
}

// OTPVerifyInput — payload verifikasi OTP
type OTPVerifyInput struct {
	Email   string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp_code" validate:"required,len=6"`
}

// RegisterResult — output dari service
type RegisterResult struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// OTPRecord — data OTP dari database
type OTPRecord struct {
	ID        int
	UserID    string
	Email     string
	Code      string
	ExpiresAt time.Time
	Used      bool
}
