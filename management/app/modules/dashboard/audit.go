package dashboard

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureAuditSchema(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
		CREATE SCHEMA IF NOT EXISTS management;
		CREATE TABLE IF NOT EXISTS management.audit_logs (
			id BIGSERIAL PRIMARY KEY,
			endpoint TEXT NOT NULL,
			method VARCHAR(16) NOT NULL,
			role VARCHAR(64),
			status_code INT NOT NULL,
			message TEXT,
			payload JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_management_audit_logs_created_at ON management.audit_logs(created_at DESC);
	`
	_, err := pool.Exec(ctx, query)
	return err
}

func AuditMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		audit := requestAudit{
			Endpoint: c.Request.URL.Path,
			Method:   c.Request.Method,
			Role:     strings.ToUpper(strings.TrimSpace(c.GetHeader("X-User-Role"))),
			Status:   c.Writer.Status(),
			Message:  strings.TrimSpace(c.GetString("audit_error")),
			Payload:  c.Request.URL.RawQuery,
		}

		if err := saveAuditLog(c.Request.Context(), pool, audit); err != nil {
			log.Printf("[MANAGEMENT] gagal menyimpan audit log: %v", err)
		}
	}
}

func setAuditError(c *gin.Context, msg string) {
	if msg == "" {
		return
	}
	c.Set("audit_error", msg)
}

func saveAuditLog(ctx context.Context, pool *pgxpool.Pool, audit requestAudit) error {
	query := `
		INSERT INTO management.audit_logs (endpoint, method, role, status_code, message, payload)
		VALUES ($1, $2, $3, $4, $5, CASE WHEN $6 = '' THEN NULL ELSE to_jsonb($6::text) END)
	`
	_, err := pool.Exec(ctx, query, audit.Endpoint, audit.Method, audit.Role, audit.Status, audit.Message, audit.Payload)
	return err
}
