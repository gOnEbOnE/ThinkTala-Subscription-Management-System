package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"notification/app/routes"
	"notification/core/database"
	"notification/core/queue"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Init(); err != nil {
		log.Fatalf("[NOTIFICATION] %v", err)
	}
	defer database.DB.Close()

	database.Migrate()
	database.Seed()

	r := gin.Default()
	// TODO(queue): Ganti _ dengan svc dan uncomment workers di bawah saat siap di-deploy
	// __ := routes.Register(r)
	svc := routes.Register(r)
	// import (
	//   "context"
	//   "os/signal"
	//   "syscall"
	//   "notification/core/queue"
	// )
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()
	// svc := routes.Register(r)
	// go svc.StartRetryWorker(ctx)
	// go queue.StartWorker(ctx, svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go svc.StartRetryWorker(ctx)
	go queue.StartWorker(ctx, svc)

	port := database.GetEnv("PORT", "5003")
	log.Printf("[NOTIFICATION] Service running on :%s", port)
	r.Run(":" + port)
}
