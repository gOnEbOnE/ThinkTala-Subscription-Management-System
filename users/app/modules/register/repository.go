package register

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

type Repository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	GetUserStatusByEmail(ctx context.Context, email string) (userID string, status string, err error)
	DeleteInactiveUser(ctx context.Context, email string) error
	CreateUser(ctx context.Context, id, name, email, password, phone, birthdate string) error
	SaveOTP(ctx context.Context, userID, email, code string, expiresAt time.Time) error
	GetValidOTP(ctx context.Context, email, code string) (*OTPRecord, error)
	MarkOTPUsed(ctx context.Context, otpID int) error
	ActivateUser(ctx context.Context, email string) error
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
}

type registerRepo struct {
	db *database.DBWrapper
}

func NewRepository(db *database.DBWrapper) Repository {
	return &registerRepo{db: db}
}

// EmailExists mengecek apakah email sudah terdaftar
func (r *registerRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM users WHERE email = $1", email,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserStatusByEmail mengambil user ID dan status berdasarkan email
func (r *registerRepo) GetUserStatusByEmail(ctx context.Context, email string) (string, string, error) {
	var userID, status string
	err := r.db.Pool.QueryRow(ctx,
		"SELECT id, status FROM users WHERE email = $1", email,
	).Scan(&userID, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", nil
		}
		return "", "", err
	}
	return userID, status, nil
}

// DeleteInactiveUser menghapus user inactive beserta OTP-nya agar bisa register ulang
func (r *registerRepo) DeleteInactiveUser(ctx context.Context, email string) error {
	// Hapus OTP dulu (foreign key)
	_, _ = r.db.Pool.Exec(ctx,
		"DELETE FROM otp_codes WHERE email = $1", email,
	)
	// Hapus user inactive
	_, err := r.db.Pool.Exec(ctx,
		"DELETE FROM users WHERE email = $1 AND status = 'inactive'", email,
	)
	return err
}

// CreateUser menyimpan user baru dengan status 'inactive'
func (r *registerRepo) CreateUser(ctx context.Context, id, name, email, password, phone, birthdate string) error {
	_, err := r.db.Pool.Exec(ctx, `
        INSERT INTO users (id, name, email, password, phone, birthdate, group_id, level_id, role_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6,
            '2e98c63f-5474-4506-826c-ded22b59b3dc',
            2,
            'df47ce1c-1455-4a20-bafe-c2b7c2ab9994',
            'inactive',
            CURRENT_TIMESTAMP
        )
    `, id, name, email, password, phone, birthdate)
	return err
}

// SaveOTP menyimpan OTP ke database
func (r *registerRepo) SaveOTP(ctx context.Context, userID, email, code string, expiresAt time.Time) error {
	// Invalidasi OTP lama yang belum digunakan
	_, _ = r.db.Pool.Exec(ctx,
		"UPDATE otp_codes SET used = TRUE WHERE email = $1 AND used = FALSE", email,
	)

	_, err := r.db.Pool.Exec(ctx, `
        INSERT INTO otp_codes (user_id, email, code, expires_at)
        VALUES ($1, $2, $3, $4)
    `, userID, email, code, expiresAt)
	return err
}

// GetValidOTP mengambil OTP yang valid (belum expired & belum digunakan)
func (r *registerRepo) GetValidOTP(ctx context.Context, email, code string) (*OTPRecord, error) {
	var otp OTPRecord
	err := r.db.Pool.QueryRow(ctx, `
        SELECT id, user_id, email, code, expires_at, used
        FROM otp_codes
        WHERE email = $1 AND code = $2 AND used = FALSE
        ORDER BY created_at DESC
        LIMIT 1
    `, email, code).Scan(&otp.ID, &otp.UserID, &otp.Email, &otp.Code, &otp.ExpiresAt, &otp.Used)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

// MarkOTPUsed menandai OTP sebagai sudah digunakan
func (r *registerRepo) MarkOTPUsed(ctx context.Context, otpID int) error {
	_, err := r.db.Pool.Exec(ctx,
		"UPDATE otp_codes SET used = TRUE WHERE id = $1", otpID,
	)
	return err
}

// ActivateUser mengubah status user menjadi 'active'
func (r *registerRepo) ActivateUser(ctx context.Context, email string) error {
	_, err := r.db.Pool.Exec(ctx,
		"UPDATE users SET status = 'active', updated_at = CURRENT_TIMESTAMP WHERE email = $1", email,
	)
	return err
}

// GetUserIDByEmail mengambil user ID berdasarkan email
func (r *registerRepo) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	var id string
	err := r.db.Pool.QueryRow(ctx,
		"SELECT id FROM users WHERE email = $1", email,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}
