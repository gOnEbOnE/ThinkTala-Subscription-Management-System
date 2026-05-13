package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"tickets/core/database"
)

type App struct {
	Port   string
	Router *http.ServeMux
	DB     *sql.DB
}

func New() *App {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2004"
	}

	db, dbErr := database.OpenDB()
	if dbErr != nil {
		log.Printf("[WARNING] Failed to connect database: %v", dbErr)
	} else {
		if err := database.MigrateSupportTicketsTable(db); err != nil {
			log.Printf("[WARNING] Failed to migrate support_tickets table: %v", err)
		}
	}

	return &App{
		Port:   port,
		Router: http.NewServeMux(),
		DB:     db,
	}
}

func (a *App) Close() {
	if a == nil || a.DB == nil {
		return
	}
	a.DB.Close()
}

func (a *App) Run() {
	fmt.Printf("Tickets service is running on port %s\n", a.Port)
	if err := http.ListenAndServe(":"+a.Port, a.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
