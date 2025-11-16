package router

import (
	"github.com/fikriahmadf/outbox-examples/internal/handlers/internal_memo"
	"github.com/gofiber/fiber/v2"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	InternalMemoHandler *internal_memo.MemoHandler
}

// Router is the router struct containing handlers.
type Router struct {
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

// ProvideDomainHandlers constructs DomainHandlers from individual handlers.
func ProvideDomainHandlers(memoHandler *internal_memo.MemoHandler) DomainHandlers {
	return DomainHandlers{
		InternalMemoHandler: memoHandler,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	if r.DomainHandlers.InternalMemoHandler != nil {
		r.DomainHandlers.InternalMemoHandler.Router(v1)
	}
}
