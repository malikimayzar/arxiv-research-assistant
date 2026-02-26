package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/malikimayzar/arxiv-research-assistant/internal/api"
	"github.com/malikimayzar/arxiv-research-assistant/internal/client"
	"github.com/malikimayzar/arxiv-research-assistant/internal/middleware"
	"github.com/malikimayzar/arxiv-research-assistant/internal/repository"
)

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "arxiv"),
		Password: getEnv("DB_PASSWORD", "arxiv_secret"),
		DBName:   getEnv("DB_NAME", "arxiv_db"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	mlClient := client.NewMLClient(getEnv("ML_SERVICE_URL", "http://localhost:8001"))
	if err := mlClient.Ping(); err != nil {
		log.Printf("ML service not reachable: %v", err)
	} else {
		log.Println("Connected to ML service")
	}

	app := fiber.New(fiber.Config{
		AppName:      "ArXiv Research Assistant v0.1.0",
		ErrorHandler: api.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} | req-id: ${respHeader:X-Request-ID}\n",
	}))

	app.Get("/health", api.HealthHandler(db, mlClient))

	log.Println("Server starting on :8080")
	log.Fatal(app.Listen(":8080"))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}