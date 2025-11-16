//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
)

// Injectors from wire.go:

func InitializeApplication() (*Application, error) {
	config := ProvideConfig()
	postgresConn := infras.ProvidePostgresConn(config)
	internalMemoRepositoryPostgres := repository.ProvideInternalMemoRepositoryPostgres(postgresConn)
	application := NewApplication(config, postgresConn, internalMemoRepositoryPostgres)
	return application, nil
}
