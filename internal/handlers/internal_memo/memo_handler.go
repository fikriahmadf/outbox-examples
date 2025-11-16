package internal_memo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
)

// CreateMemo godoc
// @Summary Create a memo
// @Description Create a new internal memo
// @Tags internal_memo
// @Accept json
// @Produce json
// @Param payload body model.CreateMemoRequest true "Memo payload"
// @Success 201 {object} model.Memo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/internal_memo/memo [post]
func (h *MemoHandler) CreateMemo(c *fiber.Ctx) error {
	var req model.CreateMemoRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn().Err(err).Msg("[MemoHandler][CreateMemo] failed to parse body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	newModel := req.ToNewModel()
	if err := h.MemoRepo.CreateMemo(c.Context(), &newModel); err != nil {
		log.Warn().Err(err).Msg("[MemoHandler][CreateMemo] failed to create memo")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create memo"})
	}

	return c.Status(fiber.StatusCreated).JSON(newModel)
}
