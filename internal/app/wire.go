//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package app

import (
	"github.com/google/wire"

	extNotifPublisherService "github.com/fikriahmadf/outbox-examples/external/domain/notif_publisher/service"
	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/internal/domain/internal_memo/repository"
	internalmemo "github.com/fikriahmadf/outbox-examples/internal/handlers/internal_memo"
	httpserver "github.com/fikriahmadf/outbox-examples/transport/http"
	"github.com/fikriahmadf/outbox-examples/transport/http/router"
)

// InitializeServer constructs the HTTP server with all handlers using compile-time DI.
func InitializeServer() (*httpserver.HTTP, error) {
	wire.Build(
		ProvideConfig,
		infras.ProvidePostgresConn,
		repository.ProvideInternalMemoRepositoryPostgres,
		wire.Bind(new(repository.InternalMemoRepository), new(*repository.InternalMemoRepositoryPostgres)),
		extNotifPublisherService.ProvideNotifPublisherService,
		wire.Bind(new(extNotifPublisherService.ExternalNotifPublisherService), new(*extNotifPublisherService.ExternalNotifPublisherServiceImpl)),
		internalmemo.ProvideMemoHandler,
		router.ProvideDomainHandlers,
		router.ProvideRouter,
		httpserver.ProvideHTTP,
	)
	return nil, nil
}
