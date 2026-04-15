package main

import (
	"log"

	"notification/app/routes"
	"notification/core/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"context"
	"os/signal"
	"syscall"
	"notification/core/queue"
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
	svc := routes.Register(r)

	// TODO(queue): Aktifkan semua baris berikut setelah workers siap di-deploy:
	// import (
	//   "context"
	//   "os/signal"
	//   "syscall"
	//   "notification/core/queue"
	// )
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// svc = routes.Register(r)
	go svc.StartRetryWorker(ctx)
	go queue.StartWorker(ctx, svc)

	port := database.GetEnv("PORT", "5003")
	log.Printf("[NOTIFICATION] Service running on :%s", port)
	r.Run(":" + port)
}
