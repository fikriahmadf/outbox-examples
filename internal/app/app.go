package app

import (
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
)

// Application centralizes the core dependencies of the service.
type Application struct {
	Config *configs.Config
	DB     *infras.PostgresConn
	InternalMemoRepo repository.InternalMemoRepository
}

// ProvideConfig exposes the singleton configuration loader to the DI container.
func ProvideConfig() *configs.Config {
	return configs.Get()
}

// NewApplication wires core dependencies into a single Application value.
func NewApplication(config *configs.Config, db *infras.PostgresConn, repo repository.InternalMemoRepository) *Application {
	log.Info().Msg("Application dependencies wired.")
	return &Application{Config: config, DB: db, InternalMemoRepo: repo}
}
