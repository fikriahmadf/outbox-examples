package app

import (
	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
)

// Application centralizes the core dependencies of the service.
type Application struct {
	Config           *configs.Config
	DB               *infras.PostgresConn
	InternalMemoRepo repository.InternalMemoRepository
}

// ProvideConfig exposes the singleton configuration loader to the DI container.
func ProvideConfig() *configs.Config {
	return configs.Get()
}
