package main

import (
	"context"
	"log"
	"time"

	"celana-api/internal/config"
	"celana-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := config.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if err := config.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	if err := config.EnsureAdminUser(ctx, db); err != nil {
		log.Fatalf("Admin seed failed: %v", err)
	}
	log.Println("Database connected, migrations applied, and admin verified")

	app := fiber.New(fiber.Config{
		AppName: "Celana (Cellix Laundry Management) API",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	routes.Setup(app, db, cfg.JWTSecret)

	log.Printf("Celana API listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
