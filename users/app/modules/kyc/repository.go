package kyc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/master-abror/zaframework/core/database"
)

// Repository mendefinisikan interface untuk akses data KYC
type Repository interface {
	NIKExists(ctx context.Context, nik string) (bool, error)
	NIKExistsForOtherUser(ctx context.Context, nik string, userID string) (bool, error)
	UserHasPendingKYC(ctx context.Context, userID string) (bool, error)
	CreateSubmission(ctx context.Context, sub *KYCSubmission) error
	GetByUserID(ctx context.Context, userID string) (*KYCSubmission, error)
	GetByID(ctx context.Context, id string) (*KYCSubmission, error)
	UpdateResubmission(ctx context.Context, id string, fullName string, nik string, address string, birthdate string, phone string, ktpImage string) error
	ResubmitKYC(ctx context.Context, userID string, fullName string, nik string, address string, birthdate string, phone string, ktpImage string) (string, error)

	// Admin methods
	ListAll(ctx context.Context, status string, search string, page int, limit int) ([]KYCListItem, int, error)
	GetDetailByID(ctx context.Context, id string) (*KYCDetailResult, error)
	UpdateStatus(ctx context.Context, id string, status string, reviewerID string, rejectReason string) error
}

type kycRepo struct {
	db *database.DBWrapper
}

// NewRepository membuat instance baru KYC Repository
func NewRepository(db *database.DBWrapper) Repository {
	return &kycRepo{db: db}
}

// NIKExists mengecek apakah NIK sudah terdaftar di tabel kyc_submissions
func (r *kycRepo) NIKExists(ctx context.Context, nik string) (bool, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM kyc_submissions WHERE nik = $1", nik,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// NIKExistsForOtherUser mengecek apakah NIK sudah terdaftar di user LAIN (bukan user yang sama)
func (r *kycRepo) NIKExistsForOtherUser(ctx context.Context, nik string, userID string) (bool, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM kyc_submissions WHERE nik = $1 AND user_id != $2", nik, userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UserHasPendingKYC mengecek apakah user sudah punya pengajuan KYC yang pending/approved
func (r *kycRepo) UserHasPendingKYC(ctx context.Context, userID string) (bool, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM kyc_submissions WHERE user_id = $1 AND status IN ('pending', 'approved')", userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateSubmission menyimpan pengajuan KYC baru ke database
func (r *kycRepo) CreateSubmission(ctx context.Context, sub *KYCSubmission) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO kyc_submissions (id, user_id, full_name, nik, address, birthdate, phone, ktp_image, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, sub.ID, sub.UserID, sub.FullName, sub.NIK, sub.Address, sub.Birthdate, sub.Phone, sub.KTPImage, sub.Status)
	return err
}

// GetByUserID mengambil data KYC terbaru berdasarkan user_id
func (r *kycRepo) GetByUserID(ctx context.Context, userID string) (*KYCSubmission, error) {
	var sub KYCSubmission
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, full_name, nik, address, birthdate::TEXT, phone, ktp_image, 
		       status, reject_reason, reviewed_by, reviewed_at, created_at, updated_at
		FROM kyc_submissions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(
		&sub.ID, &sub.UserID, &sub.FullName, &sub.NIK, &sub.Address, &sub.Birthdate,
		&sub.Phone, &sub.KTPImage, &sub.Status, &sub.RejectReason,
		&sub.ReviewedBy, &sub.ReviewedAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// GetByID mengambil data KYC berdasarkan ID submission
func (r *kycRepo) GetByID(ctx context.Context, id string) (*KYCSubmission, error) {
	var sub KYCSubmission
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, full_name, nik, address, birthdate::TEXT, phone, ktp_image,
		       status, reject_reason, reviewed_by, reviewed_at, created_at, updated_at
		FROM kyc_submissions
		WHERE id = $1
	`, id).Scan(
		&sub.ID, &sub.UserID, &sub.FullName, &sub.NIK, &sub.Address, &sub.Birthdate,
		&sub.Phone, &sub.KTPImage, &sub.Status, &sub.RejectReason,
		&sub.ReviewedBy, &sub.ReviewedAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// UpdateResubmission memperbarui KYC yang rejected: reset ke pending, hapus rejection reason
func (r *kycRepo) UpdateResubmission(ctx context.Context, id string, fullName string, nik string, address string, birthdate string, phone string, ktpImage string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE kyc_submissions
		SET full_name = $1, nik = $2, address = $3, birthdate = $4, phone = $5,
		    ktp_image = $6, status = 'pending', reject_reason = NULL,
		    reviewed_by = NULL, reviewed_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`, fullName, nik, address, birthdate, phone, ktpImage, id)
	return err
}

// ResubmitKYC — PBI-8: Transactional resubmit (check status=rejected AND update in one atomic operation)
// Returns (oldKTPImage, error). Returns specific error strings for non-rejected statuses.
func (r *kycRepo) ResubmitKYC(ctx context.Context, userID string, fullName string, nik string, address string, birthdate string, phone string, ktpImage string) (string, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback(ctx)

	// Lock the row and check current status
	var kycID, currentStatus, oldKTPImage string
	err = tx.QueryRow(ctx, `
		SELECT id, status, ktp_image
		FROM kyc_submissions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID).Scan(&kycID, &currentStatus, &oldKTPImage)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("NO_KYC_RECORD")
		}
		return "", err
	}

	if currentStatus != "rejected" {
		return "", fmt.Errorf("STATUS_NOT_REJECTED:%s", currentStatus)
	}

	// Update the record
	_, err = tx.Exec(ctx, `
		UPDATE kyc_submissions
		SET full_name = $1, nik = $2, address = $3, birthdate = $4, phone = $5,
		    ktp_image = $6, status = 'pending', reject_reason = NULL,
		    reviewed_by = NULL, reviewed_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`, fullName, nik, address, birthdate, phone, ktpImage, kycID)
	if err != nil {
		return "", fmt.Errorf("gagal memperbarui data KYC: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("gagal commit transaksi: %v", err)
	}

	return oldKTPImage, nil
}

// ========== ADMIN METHODS ==========

// ListAll mengambil daftar semua pengajuan KYC dengan filter dan pagination
func (r *kycRepo) ListAll(ctx context.Context, status string, search string, page int, limit int) ([]KYCListItem, int, error) {
	countQuery := `SELECT COUNT(*) FROM kyc_submissions ks LEFT JOIN users u ON ks.user_id::text = u.id::text WHERE 1=1`
	dataQuery := `SELECT ks.id, ks.user_id, ks.full_name, ks.nik, COALESCE(u.email, '') as email, ks.status, ks.created_at 
		FROM kyc_submissions ks LEFT JOIN users u ON ks.user_id::text = u.id::text WHERE 1=1`

	args := []any{}
	argIdx := 1

	if status != "" && status != "all" {
		countQuery += fmt.Sprintf(" AND ks.status = $%d", argIdx)
		dataQuery += fmt.Sprintf(" AND ks.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if search != "" {
		countQuery += fmt.Sprintf(" AND (ks.full_name ILIKE $%d OR ks.nik ILIKE $%d OR u.email ILIKE $%d)", argIdx, argIdx, argIdx)
		dataQuery += fmt.Sprintf(" AND (ks.full_name ILIKE $%d OR ks.nik ILIKE $%d OR u.email ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	var total int
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery += " ORDER BY ks.created_at DESC"
	offset := (page - 1) * limit
	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []KYCListItem
	for rows.Next() {
		var item KYCListItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.FullName, &item.NIK, &item.Email, &item.Status, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

// GetDetailByID mengambil detail lengkap KYC beserta email user
func (r *kycRepo) GetDetailByID(ctx context.Context, id string) (*KYCDetailResult, error) {
	var d KYCDetailResult
	err := r.db.Pool.QueryRow(ctx, `
		SELECT ks.id, ks.user_id, ks.full_name, ks.nik, ks.address, ks.birthdate::TEXT, ks.phone, ks.ktp_image,
		       COALESCE(u.email, '') as email, ks.status, ks.reject_reason, ks.reviewed_by, ks.reviewed_at,
		       ks.created_at, ks.updated_at
		FROM kyc_submissions ks
		LEFT JOIN users u ON ks.user_id::text = u.id::text
		WHERE ks.id = $1
	`, id).Scan(
		&d.ID, &d.UserID, &d.FullName, &d.NIK, &d.Address, &d.Birthdate,
		&d.Phone, &d.KTPImage, &d.Email, &d.Status, &d.RejectReason,
		&d.ReviewedBy, &d.ReviewedAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

// UpdateStatus mengubah status KYC (approve/reject)
func (r *kycRepo) UpdateStatus(ctx context.Context, id string, status string, reviewerID string, rejectReason string) error {
	var err error
	if status == "rejected" {
		_, err = r.db.Pool.Exec(ctx, `
			UPDATE kyc_submissions 
			SET status = $1, reviewed_by = $2, reject_reason = $3, reviewed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE id = $4
		`, status, reviewerID, rejectReason, id)
	} else {
		_, err = r.db.Pool.Exec(ctx, `
			UPDATE kyc_submissions 
			SET status = $1, reviewed_by = $2, reject_reason = NULL, reviewed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
		`, status, reviewerID, id)
	}
	return err
}
