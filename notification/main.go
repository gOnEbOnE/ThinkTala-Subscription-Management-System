package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	log.Println("[NOTIFICATION] Connected to database")

	r := gin.Default()

	// ===== Compliance Dashboard Pages =====
	r.GET("/compliance/dashboard", func(c *gin.Context) {
		c.File("./public/views/dashboard.html")
	})
	r.GET("/compliance/kyc-review", func(c *gin.Context) {
		c.File("./public/views/kyc-review.html")
	})

	// ===== Existing Notification API =====
	api := r.Group("/api/notifications")
	{
		api.GET("", listNotifications)
		api.GET("/:id", getNotification)
		api.POST("", createNotification)
		api.PUT("/:id", updateNotification)
		api.DELETE("/:id", deleteNotification)
		api.GET("/public", listPublicNotifications)
	}

	port := getEnv("PORT", "5003")
	log.Printf("[NOTIFICATION] Service running on :%s", port)
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type Notification struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`
	TargetRole string    `json:"target_role"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  *string   `json:"created_by,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func listNotifications(c *gin.Context) {
	rows, err := db.Query(context.Background(),
		`SELECT id, title, message, type, target_role, is_active, created_at, created_by, updated_at 
         FROM notifications ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
		var n Notification
		rows.Scan(&n.ID, &n.Title, &n.Message, &n.Type, &n.TargetRole, &n.IsActive, &n.CreatedAt, &n.CreatedBy, &n.UpdatedAt)
		list = append(list, n)
	}
	if list == nil {
		list = []Notification{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func listPublicNotifications(c *gin.Context) {
	role := c.Query("role")
	rows, err := db.Query(context.Background(),
		`SELECT id, title, message, type, target_role, created_at 
         FROM notifications 
         WHERE is_active = TRUE AND (target_role = 'all' OR target_role = $1)
         ORDER BY created_at DESC LIMIT 20`, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []map[string]any
	for rows.Next() {
		var id, title, msg, typ, target string
		var createdAt time.Time
		rows.Scan(&id, &title, &msg, &typ, &target, &createdAt)
		list = append(list, map[string]any{
			"id": id, "title": title, "message": msg, "type": typ,
			"target_role": target, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]any{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func getNotification(c *gin.Context) {
	id := c.Param("id")
	var n Notification
	err := db.QueryRow(context.Background(),
		`SELECT id, title, message, type, target_role, is_active, created_at, created_by, updated_at 
         FROM notifications WHERE id = $1`, id,
	).Scan(&n.ID, &n.Title, &n.Message, &n.Type, &n.TargetRole, &n.IsActive, &n.CreatedAt, &n.CreatedBy, &n.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": n})
}

func createNotification(c *gin.Context) {
	var body struct {
		Title      string `json:"title" binding:"required"`
		Message    string `json:"message" binding:"required"`
		Type       string `json:"type"`
		TargetRole string `json:"target_role"`
		IsActive   *bool  `json:"is_active"`
		CreatedBy  string `json:"created_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	if body.Type == "" {
		body.Type = "info"
	}
	if body.TargetRole == "" {
		body.TargetRole = "all"
	}
	isActive := true
	if body.IsActive != nil {
		isActive = *body.IsActive
	}

	var createdBy *string
	if body.CreatedBy != "" {
		createdBy = &body.CreatedBy
	}

	_, err := db.Exec(context.Background(),
		`INSERT INTO notifications (id, title, message, type, target_role, is_active, created_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`, id, body.Title, body.Message, body.Type, body.TargetRole, isActive, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}, "message": "Created"})
}

func updateNotification(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Title      string `json:"title"`
		Message    string `json:"message"`
		Type       string `json:"type"`
		TargetRole string `json:"target_role"`
		IsActive   *bool  `json:"is_active"`
		UpdatedBy  string `json:"updated_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedBy *string
	if body.UpdatedBy != "" {
		updatedBy = &body.UpdatedBy
	}

	_, err := db.Exec(context.Background(),
		`UPDATE notifications SET title=COALESCE(NULLIF($2,''),title), message=COALESCE(NULLIF($3,''),message),
         type=COALESCE(NULLIF($4,''),type), target_role=COALESCE(NULLIF($5,''),target_role),
         is_active=COALESCE($6, is_active), updated_by=$7, updated_at=CURRENT_TIMESTAMP
         WHERE id=$1`,
		id, body.Title, body.Message, body.Type, body.TargetRole, body.IsActive, updatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func deleteNotification(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec(context.Background(), "DELETE FROM notifications WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// Helper for JSON marshal (unused but kept for reference)
func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
