package service

import (
	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/external/domain/notif_publisher/service"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
)

type InternalMemoService interface {
	EmailOutboxService
}

type InternalMemoServiceImpl struct {
	Config                        *configs.Config
	InternalMemoRepository        repository.InternalMemoRepository
	ExternalNotifPublisherService service.ExternalNotifPublisherService
}

func ProvideInternalMemoService(cfg *configs.Config, internalMemoRepository repository.InternalMemoRepository, externalNotifPublisherService service.ExternalNotifPublisherService) *InternalMemoServiceImpl {
	return &InternalMemoServiceImpl{
		Config:                        cfg,
		InternalMemoRepository:        internalMemoRepository,
		ExternalNotifPublisherService: externalNotifPublisherService,
	}
}
