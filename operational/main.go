package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func main() {
	godotenv.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "thinknalyze"),
		getEnv("DB_SSLMODE", "disable"),
	)

	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Gagal koneksi DB: %v", err)
	}
	defer db.Close()
	log.Println("[OPERATIONAL] Connected to database")

	r := gin.Default()

	// Serve static assets
	r.Static("/ops/assets", "./public/assets")

	// Dashboard pages
	r.GET("/ops/dashboard", func(c *gin.Context) {
		c.File("./public/views/dashboard.html")
	})
	r.GET("/ops/notifications", func(c *gin.Context) {
		c.File("./public/views/notifications.html")
	})
	r.GET("/ops/notification-templates", func(c *gin.Context) {
        c.File("./public/views/notification-templates.html")
    })
	r.GET("/ops/subscriptions", func(c *gin.Context) {
		c.File("./public/views/subscriptions.html")
	})
	r.GET("/ops/users", func(c *gin.Context) {
		c.File("./public/views/users.html")
	})

	// API routes (existing + new)
	api := r.Group("/api/operational")
	{
		api.GET("/stats", getDashboardStats)
	}

	// KYC API (if exists)
	kyc := r.Group("/api/kyc")
	{
		kyc.GET("", listKYC)
		kyc.PUT("/:id/review", reviewKYC)
	}

	port := getEnv("PORT", "5005")
	log.Printf("[OPERATIONAL] Service running on :%s", port)
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDashboardStats(c *gin.Context) {
	var totalUsers, totalNotif, totalSubs, pendingKYC int

	db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&totalUsers)
	db.QueryRow(context.Background(), "SELECT COUNT(*) FROM notifications WHERE is_active=TRUE").Scan(&totalNotif)
	db.QueryRow(context.Background(), "SELECT COUNT(*) FROM subscription_packages WHERE is_active=TRUE").Scan(&totalSubs)
	db.QueryRow(context.Background(), "SELECT COUNT(*) FROM kyc_documents WHERE status='pending'").Scan(&pendingKYC)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_users":          totalUsers,
			"active_notifications": totalNotif,
			"subscription_plans":   totalSubs,
			"pending_kyc":          pendingKYC,
		},
	})
}

func listKYC(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	var query string
	var args []any

	if status != "" {
		query = `SELECT k.id, k.user_id, u.name, u.email, k.document_type, k.file_url, k.status, k.reject_reason, k.created_at
                 FROM kyc_documents k JOIN users u ON u.id = k.user_id WHERE k.status = $1 ORDER BY k.created_at DESC`
		args = append(args, status)
	} else {
		query = `SELECT k.id, k.user_id, u.name, u.email, k.document_type, k.file_url, k.status, k.reject_reason, k.created_at
                 FROM kyc_documents k JOIN users u ON u.id = k.user_id ORDER BY k.created_at DESC`
	}

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []map[string]any
	for rows.Next() {
		var id, userID, name, email, docType, fileURL, st string
		var rejectReason *string
		var createdAt time.Time
		rows.Scan(&id, &userID, &name, &email, &docType, &fileURL, &st, &rejectReason, &createdAt)
		list = append(list, map[string]any{
			"id": id, "user_id": userID, "name": name, "email": email,
			"document_type": docType, "file_url": fileURL, "status": st,
			"reject_reason": rejectReason, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]any{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func reviewKYC(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status       string `json:"status" binding:"required"`
		RejectReason string `json:"reject_reason"`
		ReviewedBy   string `json:"reviewed_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(context.Background(),
		`UPDATE kyc_documents SET status=$2, reject_reason=$3, reviewed_by=$4, reviewed_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1`,
		id, body.Status, body.RejectReason, body.ReviewedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "KYC reviewed"})
}
