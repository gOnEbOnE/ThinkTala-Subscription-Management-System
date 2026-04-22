package database

import (
	"context"
	"database/sql"
	"time"
)

func MigrateSupportTicketsTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS support_tickets (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			reporter_name VARCHAR(255) NOT NULL DEFAULT '',
			reporter_email VARCHAR(255) NOT NULL DEFAULT '',
			title TEXT NOT NULL,
			category VARCHAR(50) NOT NULL DEFAULT 'MASALAH_TEKNIS',
			description TEXT NOT NULL DEFAULT '',
			status VARCHAR(20) NOT NULL DEFAULT 'ON PROCESS',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS support_ticket_replies (
			id UUID PRIMARY KEY,
			ticket_id UUID NOT NULL,
			admin_id UUID NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT support_ticket_replies_ticket_id_fkey
				FOREIGN KEY (ticket_id) REFERENCES support_tickets(id) ON DELETE CASCADE
		);

		ALTER TABLE support_tickets DROP CONSTRAINT IF EXISTS support_tickets_user_id_fkey;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS reporter_name VARCHAR(255) NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS reporter_email VARCHAR(255) NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS category VARCHAR(50) NOT NULL DEFAULT 'MASALAH_TEKNIS';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_name VARCHAR(255);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_mime VARCHAR(100);
		ALTER TABLE support_tickets ADD COLUMN IF NOT EXISTS attachment_data BYTEA;
		ALTER TABLE support_ticket_replies ADD COLUMN IF NOT EXISTS ticket_id UUID;
		ALTER TABLE support_ticket_replies ADD COLUMN IF NOT EXISTS admin_id UUID;
		ALTER TABLE support_ticket_replies ADD COLUMN IF NOT EXISTS message TEXT;
		ALTER TABLE support_ticket_replies ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

		ALTER TABLE support_ticket_replies DROP CONSTRAINT IF EXISTS support_ticket_replies_ticket_id_fkey;
		ALTER TABLE support_ticket_replies
			ADD CONSTRAINT support_ticket_replies_ticket_id_fkey
			FOREIGN KEY (ticket_id) REFERENCES support_tickets(id) ON DELETE CASCADE;

		CREATE INDEX IF NOT EXISTS idx_support_tickets_created_at ON support_tickets(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_support_tickets_user_id ON support_tickets(user_id);
		CREATE INDEX IF NOT EXISTS idx_support_tickets_status ON support_tickets(status);
		CREATE INDEX IF NOT EXISTS idx_support_ticket_replies_ticket_id ON support_ticket_replies(ticket_id);
		CREATE INDEX IF NOT EXISTS idx_support_ticket_replies_created_at ON support_ticket_replies(created_at ASC);
	`)

	return err
}
