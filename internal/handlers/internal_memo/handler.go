package internal_memo

import (
	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
	"github.com/gofiber/fiber/v2"
)

type MemoHandler struct {
	Config                 *configs.Config
	InternalMemoRepository repository.InternalMemoRepository
}

func ProvideMemoHandler(config *configs.Config, intMemoRepo repository.InternalMemoRepository) *MemoHandler {
	return &MemoHandler{
		Config:                 config,
		InternalMemoRepository: intMemoRepo,
	}
}

func (h *MemoHandler) Router(r fiber.Router) {
	r.Route("/internal_memo", func(router fiber.Router) {
		router.Route("/memo", func(router fiber.Router) {
			router.Post("", h.CreateMemo)
		})
	})
}
