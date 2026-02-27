package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/malikimayzar/arxiv-research-assistant/internal/client"
	"github.com/malikimayzar/arxiv-research-assistant/internal/repository"
)

type QueryRequest struct {
	Query   string `json:"query"`
	TopK    int    `json:"top_k"`
	Model   string `json:"model"`
	ArxivID string `json:"arxiv_id,omitempty"`
}

type QueryResponse struct {
	Answer       string          `json:"answer"`
	Sources      []client.Source `json:"sources"`
	RetrievalMs  int             `json:"retrieval_ms"`
	GenerationMs int             `json:"generation_ms"`
	Model        string          `json:"model"`
	LogID        string          `json:"log_id,omitempty"`
}

func QueryHandler(ml *client.MLClient, logRepo repository.QueryLogRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req QueryRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if req.Query == "" {
			return fiber.NewError(fiber.StatusBadRequest, "query cannot be empty")
		}
		if req.TopK == 0 {
			req.TopK = 5
		}
		if req.Model == "" {
			req.Model = "phi3:mini"
		}

		result, err := ml.Query(req.Query, req.TopK, req.Model, req.ArxivID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// Log ke PostgreSQL — sync, tapi dengan timeout terpisah
		var logID string
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logged, err := logRepo.Create(ctx, repository.CreateQueryLogParams{
				Query:        req.Query,
				Answer:       result.Answer,
				RetrievalMs:  result.RetrievalMs,
				GenerationMs: result.GenerationMs,
				Model:        result.Model,
			})
			if err != nil {
				fmt.Printf("⚠️  query log failed: %v\n", err)
				return
			}
			_ = logged
		}()

		return c.JSON(QueryResponse{
			Answer:       result.Answer,
			Sources:      result.Sources,
			RetrievalMs:  result.RetrievalMs,
			GenerationMs: result.GenerationMs,
			Model:        result.Model,
			LogID:        logID,
		})
	}
}

func QueryHistoryHandler(logRepo repository.QueryLogRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit := c.QueryInt("limit", 20)

		logs, err := logRepo.List(c.Context(), limit)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch history")
		}

		return c.JSON(fiber.Map{
			"logs":  logs,
			"total": len(logs),
		})
	}
}
