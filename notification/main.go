package main

import (
	"log"

	"notification/app/routes"
	"notification/core/database"

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
	routes.Register(r)

	port := database.GetEnv("PORT", "5003")
	log.Printf("[NOTIFICATION] Service running on :%s", port)
	r.Run(":" + port)
}
