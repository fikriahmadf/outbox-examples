package internal_memo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/model"
)

func (h *MemoHandler) CreateMemo(c *fiber.Ctx) error {
	var req model.Memo
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
