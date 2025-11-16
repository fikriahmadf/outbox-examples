package internal_memo

import (
	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
	"github.com/gofiber/fiber/v2"
)

type MemoHandler struct {
	Config   *configs.Config
	MemoRepo repository.MemoRepository
}

func ProvideMemoHandler(config *configs.Config, memoRepo repository.MemoRepository) *MemoHandler {
	return &MemoHandler{
		Config:   config,
		MemoRepo: memoRepo,
	}
}

func (h *MemoHandler) Router(r fiber.Router) {
    r.Route("/internal_memo", func(router fiber.Router) {
        router.Route("/memo", func(router fiber.Router) {
            router.Post("/", h.CreateMemo)
        })
    })
}
