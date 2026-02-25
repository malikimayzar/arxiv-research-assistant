package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/malikimayzar/arxiv-research-assistant/internal/client"
)

type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Services map[string]string `json:"services"`
}

func HealthHandler(db *sqlx.DB, ml *client.MLClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		services := map[string]string{
			"postgres": "ok",
			"qdrant":   "pending",
			"ml":       "ok",
		}

		overallStatus := "ok"

		if db == nil || db.Ping() != nil {
			services["postgres"] = "error"
			overallStatus = "degraded"
		}

		if ml == nil || ml.Ping() != nil {
			services["ml"] = "error"
			overallStatus = "degraded"
		}

		return c.JSON(HealthResponse{
			Status:   overallStatus,
			Version:  "0.1.0",
			Services: services,
		})
	}
}
