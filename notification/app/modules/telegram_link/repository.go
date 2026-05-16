package telegramlink

import (
	"context"
	"fmt"
	"strings"

	"notification/core/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository menyediakan akses query ke kolom telegram_chat_id di tabel users.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository membuat instance Repository baru.
func NewRepository() *Repository {
	return &Repository{db: database.DB}
}

// SaveChatID menyimpan telegram_chat_id untuk user tertentu.
// Jika kolom belum ada, query akan gagal dan mengembalikan error deskriptif.
func (r *Repository) SaveChatID(userID, chatID string) error {
	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("user_id wajib diisi")
	}
	if strings.TrimSpace(chatID) == "" {
		return fmt.Errorf("chat_id wajib diisi")
	}
	tag, err := r.db.Exec(context.Background(), `
		UPDATE users
		SET telegram_chat_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, userID, chatID)
	if err != nil {
		return fmt.Errorf("gagal menyimpan chat_id: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}
	return nil
}

// ClearChatID menghapus telegram_chat_id (unlink) untuk user tertentu.
func (r *Repository) ClearChatID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("user_id wajib diisi")
	}
	tag, err := r.db.Exec(context.Background(), `
		UPDATE users
		SET telegram_chat_id = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("gagal menghapus chat_id: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}
	return nil
}

// GetChatID mengambil telegram_chat_id user saat ini.
func (r *Repository) GetChatID(userID string) (string, error) {
	var chatID *string
	err := r.db.QueryRow(context.Background(), `
		SELECT telegram_chat_id FROM users WHERE id = $1
	`, userID).Scan(&chatID)
	if err != nil {
		return "", fmt.Errorf("user tidak ditemukan")
	}
	if chatID == nil {
		return "", nil
	}
	return *chatID, nil
}
