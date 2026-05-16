package main

import (
	"context"
	"log"

	"management/app/modules/dashboard"
	"management/app/routes"
	"management/core/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load shared DB creds from users/.env first, then local management/.env (service port must not inherit users' port=).
	for _, envPath := range []string{"../users/.env", ".env", "../.env"} {
		_ = godotenv.Load(envPath)
	}

	pool, err := database.NewPoolFromEnvAllowNil(context.Background())
	if err != nil {
		log.Fatalf("[MANAGEMENT] gagal koneksi database: %v", err)
	}

	if pool == nil {
		r := gin.Default()
		r.GET("/health", func(c *gin.Context) {
			c.JSON(503, gin.H{
				"status":  "degraded",
				"service": "management",
				"detail":  "PostgreSQL tidak tersedia. Jalankan `docker compose up -d postgres` atau set kredensial di users/.env; mode ini aktif jika POSTGRES_OPTIONAL=true.",
			})
		})
		port := database.GetServicePort("MANAGEMENT_PORT", "5006")
		log.Printf("[MANAGEMENT] degraded mode (no DB) on :%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("[MANAGEMENT] gagal menjalankan service: %v", err)
		}
		return
	}

	defer pool.Close()

	if err := dashboard.EnsureAuditSchema(context.Background(), pool); err != nil {
		log.Fatalf("[MANAGEMENT] gagal menyiapkan schema audit: %v", err)
	}

	repo := dashboard.NewRepository(pool)
	service := dashboard.NewService(repo)
	handler := dashboard.NewHandler(service)

	r := gin.Default()
	r.Use(dashboard.AuditMiddleware(pool))
	routes.Register(r, handler)

	port := database.GetServicePort("MANAGEMENT_PORT", "5006")
	log.Printf("[MANAGEMENT] service running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[MANAGEMENT] gagal menjalankan service: %v", err)
	}
}
