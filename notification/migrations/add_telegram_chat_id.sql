-- Migration: Tambah kolom telegram_chat_id ke tabel users
-- Dijalankan sekali di database production/staging.
-- Kolom ini menyimpan Telegram Chat ID user untuk notifikasi personal.

ALTER TABLE users
ADD COLUMN IF NOT EXISTS telegram_chat_id VARCHAR(50) DEFAULT NULL;

-- Index untuk mempercepat query GetTelegramRecipients yang filter by role
-- dan check telegram_chat_id IS NOT NULL
CREATE INDEX IF NOT EXISTS idx_users_telegram_chat_id
ON users (telegram_chat_id)
WHERE telegram_chat_id IS NOT NULL;

-- Verifikasi
SELECT COUNT(*) AS total_users, COUNT(telegram_chat_id) AS users_with_telegram
FROM users;
