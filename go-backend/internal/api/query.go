package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/malikimayzar/arxiv-research-assistant/internal/client"
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
}

func QueryHandler(ml *client.MLClient) fiber.Handler {
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

		return c.JSON(QueryResponse{
			Answer:       result.Answer,
			Sources:      result.Sources,
			RetrievalMs:  result.RetrievalMs,
			GenerationMs: result.GenerationMs,
			Model:        result.Model,
		})
	}
}
