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
	log.Println("[SUBSCRIPTION] Connected to database")

	r := gin.Default()

	// ===== Client Dashboard Pages =====
	r.GET("/client/dashboard", func(c *gin.Context) {
		c.File("./public/views/dashboard.html")
	})
	r.GET("/client/subscription", func(c *gin.Context) {
		c.File("./public/views/subscription.html")
	})
	r.GET("/client/kyc", func(c *gin.Context) {
		c.File("./public/views/kyc.html")
	})

	// ===== Existing Subscription API =====
	api := r.Group("/api/subscriptions")
	{
		api.GET("", listPackages)
		api.GET("/:id", getPackage)
		api.POST("", createPackage)
		api.PUT("/:id", updatePackage)
		api.DELETE("/:id", deletePackage)
	}

	port := getEnv("PORT", "5004")
	log.Printf("[SUBSCRIPTION] Service running on :%s", port)
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type Package struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	DurationDays int       `json:"duration_days"`
	Features     []string  `json:"features"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func scanPackage(scan func(dest ...any) error) (Package, error) {
	var p Package
	var featuresJSON []byte
	err := scan(&p.ID, &p.Name, &p.Code, &p.Description, &p.Price, &p.DurationDays, &featuresJSON, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, err
	}
	json.Unmarshal(featuresJSON, &p.Features)
	if p.Features == nil {
		p.Features = []string{}
	}
	return p, nil
}

func listPackages(c *gin.Context) {
	rows, err := db.Query(context.Background(),
		`SELECT id, name, code, description, price, duration_days, features, is_active, created_at, updated_at
         FROM subscription_packages ORDER BY price ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []Package
	for rows.Next() {
		p, err := scanPackage(rows.Scan)
		if err != nil {
			continue
		}
		list = append(list, p)
	}
	if list == nil {
		list = []Package{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func getPackage(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow(context.Background(),
		`SELECT id, name, code, description, price, duration_days, features, is_active, created_at, updated_at
         FROM subscription_packages WHERE id = $1`, id)
	p, err := scanPackage(row.Scan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func createPackage(c *gin.Context) {
	var body struct {
		Name         string   `json:"name" binding:"required"`
		Code         string   `json:"code" binding:"required"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		DurationDays int      `json:"duration_days"`
		Features     []string `json:"features"`
		IsActive     *bool    `json:"is_active"`
		CreatedBy    string   `json:"created_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	if body.DurationDays == 0 {
		body.DurationDays = 30
	}
	isActive := true
	if body.IsActive != nil {
		isActive = *body.IsActive
	}
	if body.Features == nil {
		body.Features = []string{}
	}

	featuresJSON, _ := json.Marshal(body.Features)

	var createdBy *string
	if body.CreatedBy != "" {
		createdBy = &body.CreatedBy
	}

	_, err := db.Exec(context.Background(),
		`INSERT INTO subscription_packages (id, name, code, description, price, duration_days, features, is_active, created_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, id, body.Name, body.Code, body.Description, body.Price, body.DurationDays, featuresJSON, isActive, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}, "message": "Created"})
}

func updatePackage(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Name         string   `json:"name"`
		Code         string   `json:"code"`
		Description  string   `json:"description"`
		Price        *float64 `json:"price"`
		DurationDays *int     `json:"duration_days"`
		Features     []string `json:"features"`
		IsActive     *bool    `json:"is_active"`
		UpdatedBy    string   `json:"updated_by"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	featuresJSON, _ := json.Marshal(body.Features)

	var updatedBy *string
	if body.UpdatedBy != "" {
		updatedBy = &body.UpdatedBy
	}

	_, err := db.Exec(context.Background(),
		`UPDATE subscription_packages SET 
         name=COALESCE(NULLIF($2,''),name), code=COALESCE(NULLIF($3,''),code),
         description=COALESCE(NULLIF($4,''),description),
         price=COALESCE($5, price), duration_days=COALESCE($6, duration_days),
         features=COALESCE($7, features), is_active=COALESCE($8, is_active),
         updated_by=$9, updated_at=CURRENT_TIMESTAMP
         WHERE id=$1`,
		id, body.Name, body.Code, body.Description, body.Price, body.DurationDays, featuresJSON, body.IsActive, updatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func deletePackage(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec(context.Background(), "DELETE FROM subscription_packages WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
