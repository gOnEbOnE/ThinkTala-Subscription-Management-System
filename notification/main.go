package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"notification/app/routes"
	"notification/core/database"
	"notification/core/queue"
	"notification/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := database.Init(); err != nil {
		log.Fatalf("[NOTIFICATION] %v", err)
	}
	defer database.DB.Close()

	database.Migrate()
	database.Seed()

	r := gin.Default()
	svc := routes.Register(r)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go svc.StartRetryWorker(ctx)
	go queue.StartWorker(ctx, svc)

	// Mulai Telegram Bot Listener
	utils.StartTelegramListener(ctx)

	port := database.GetEnv("PORT", "5003")
	log.Printf("[NOTIFICATION] Service running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[NOTIFICATION] Server failed: %v", err)
	}
}
