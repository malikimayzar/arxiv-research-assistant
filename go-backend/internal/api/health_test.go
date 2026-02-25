package api_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/malikimayzar/arxiv-research-assistant/internal/api"
)

func TestHealthHandler_NoDB(t *testing.T) {
	app := fiber.New()
	app.Get("/health", api.HealthHandler(nil))

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result api.HealthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if result.Services["postgres"] != "error" {
		t.Errorf("expected postgres 'error' when db is nil, got '%s'", result.Services["postgres"])
	}
}
