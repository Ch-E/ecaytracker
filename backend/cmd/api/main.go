package main

import (
	"log"

	"ecaytracker/backend/config"
	"ecaytracker/backend/internal/api"
	appdb "ecaytracker/backend/internal/db"
)

func main() {
	cfg := config.Load()

	pool, err := appdb.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Printf("Database connected. Starting API on :%s (env=%s)", cfg.Port, cfg.Env)

	router := api.NewRouter(pool, cfg.FrontendURL)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
