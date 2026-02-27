package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/malikimayzar/arxiv-research-assistant/internal/repository"
)

type PapersHandler struct {
	repo repository.PaperRepository
}

func NewPapersHandler(repo repository.PaperRepository) *PapersHandler {
	return &PapersHandler{repo: repo}
}

func (h *PapersHandler) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	papers, total, err := h.repo.List(c.Context(), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch papers")
	}

	return c.JSON(fiber.Map{
		"papers": papers,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

func (h *PapersHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid paper ID")
	}

	paper, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "paper not found")
	}

	return c.JSON(paper)
}

func (h *PapersHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid paper ID")
	}

	if err := h.repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete paper")
	}

	return c.JSON(fiber.Map{"deleted": true, "id": id})
}
