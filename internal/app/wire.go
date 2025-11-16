//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package app

import (
	"github.com/google/wire"

	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
)

// InitializeApplication constructs an Application using compile-time dependency injection.
func InitializeApplication() (*Application, error) {
	wire.Build(
		ProvideConfig,
		infras.ProvidePostgresConn,
		repository.ProvideInternalMemoRepositoryPostgres,
		wire.Bind(new(repository.InternalMemoRepository), new(*repository.InternalMemoRepositoryPostgres)),
		NewApplication,
	)
	return nil, nil
}
