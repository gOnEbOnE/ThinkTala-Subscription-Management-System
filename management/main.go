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
	for _, envPath := range []string{".env", "../users/.env", "../.env"} {
		_ = godotenv.Load(envPath)
	}

	pool, err := database.NewPoolFromEnv(context.Background())
	if err != nil {
		log.Fatalf("[MANAGEMENT] gagal koneksi database: %v", err)
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
