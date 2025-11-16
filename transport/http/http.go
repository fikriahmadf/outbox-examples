package http

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/infras"
	"github.com/fikriahmadf/outbox-examples/shared/logger"
	"github.com/fikriahmadf/outbox-examples/transport/http/router"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ServerState is an indicator if this server's state.
type ServerState int

const (
	// ServerStateReady indicates that the server is ready to serve.
	ServerStateReady ServerState = iota + 1
	// ServerStateInGracePeriod indicates that the server is in its grace
	// period and will shut down after it is done cleaning up.
	ServerStateInGracePeriod
	// ServerStateInCleanupPeriod indicates that the server no longer
	// responds to any requests, is cleaning up its internal state, and
	// will shut down shortly.
	ServerStateInCleanupPeriod
)

// HTTP is the HTTP server.
type HTTP struct {
	Config *configs.Config
	DB     *infras.PostgresConn
	Router router.Router
	State  ServerState
	fiber  *fiber.App
}

// ProvideHTTP is the provider for HTTP.
func ProvideHTTP(db *infras.PostgresConn, config *configs.Config, router router.Router) *HTTP {
	return &HTTP{
		DB:     db,
		Config: config,
		Router: router,
	}
}

// SetupAndServe sets up the server and gets it up and running.
func (h *HTTP) SetupAndServe() {
	h.fiber = fiber.New()
	h.setupSwaggerDocs()
	h.setupRoutes()
	h.setupGracefulShutdown()
	h.State = ServerStateReady

	h.logServerInfo()

	log.Info().Str("port", h.Config.Server.Port).Msg("Starting up HTTP server.")

	err := h.fiber.Listen(":" + h.Config.Server.Port)
	if err != nil {
		logger.ErrorWithStack(err)
	}
}

func (h *HTTP) setupSwaggerDocs() {
	// TODO: integrate fiber-swagger if needed
}

func (h *HTTP) setupRoutes() {
	h.fiber.Get("/health", h.HealthCheck)
	h.Router.SetupRoutes(h.fiber)
}

func (h *HTTP) setupGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go h.respondToSigterm(done)
}

func (h *HTTP) respondToSigterm(done chan os.Signal) {
	<-done
	defer os.Exit(0)

	shutdownConfig := h.Config.Server.Shutdown

	log.Info().Msg("Received SIGTERM.")
	log.Info().Int64("seconds", shutdownConfig.GracePeriodSeconds).Msg("Entering grace period.")
	h.State = ServerStateInGracePeriod
	time.Sleep(time.Duration(shutdownConfig.GracePeriodSeconds) * time.Second)

	log.Info().Int64("seconds", shutdownConfig.CleanupPeriodSeconds).Msg("Entering cleanup period.")
	h.State = ServerStateInCleanupPeriod
	time.Sleep(time.Duration(shutdownConfig.CleanupPeriodSeconds) * time.Second)

	log.Info().Msg("Cleaning up completed. Shutting down now.")
}

func (h *HTTP) logServerInfo() {
	h.logCORSConfigInfo()
}

func (h *HTTP) logCORSConfigInfo() {
	corsConfig := h.Config.App.CORS
	corsHeaderInfo := "CORS Header"
	if corsConfig.Enable {
		log.Info().Msg("CORS Headers and Handlers are enabled.")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Credentials: %t", corsConfig.AllowCredentials)).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Headers: %s", strings.Join(corsConfig.AllowedHeaders, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Methods: %s", strings.Join(corsConfig.AllowedMethods, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Allow-Origin: %s", strings.Join(corsConfig.AllowedOrigins, ", "))).Msg("")
		log.Info().Str(corsHeaderInfo, fmt.Sprintf("Access-Control-Max-Age: %d", corsConfig.MaxAgeSeconds)).Msg("")
	} else {
		log.Info().Msg("CORS Headers are disabled.")
	}
}

func (h *HTTP) serverStateMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch h.State {
		case ServerStateReady:
			// Server is ready to serve, don't do anything.
			return next(c)
		case ServerStateInGracePeriod:
			// Server is in grace period. Issue a warning message and continue
			// serving as usual.
			log.Warn().Msg("SERVER IS IN GRACE PERIOD")
			return next(c)
		case ServerStateInCleanupPeriod:
			// Server is in cleanup period. Stop the request from actually
			// invoking any domain services and respond appropriately.
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "SERVER IS SHUTTING DOWN"})
		}
		return nil
	}
}

func (h *HTTP) setupCORS() {
	corsConfig := h.Config.App.CORS
	if corsConfig.Enable {
		h.fiber.Use(func(c *fiber.Ctx) error {
			c.Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", corsConfig.AllowCredentials))
			c.Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
			c.Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
			c.Set("Access-Control-Allow-Origin", strings.Join(corsConfig.AllowedOrigins, ", "))
			c.Set("Access-Control-Max-Age", fmt.Sprintf("%d", corsConfig.MaxAgeSeconds))
			return c.Next()
		})
	}
}

// HealthCheck performs a health check on the server. Usually required by
// Kubernetes to check if the service is healthy.
// @Summary Health Check
// @Description Health Check Endpoint
// @Tags service
// @Produce json
// @Accept json
// @Success 200 {object} response.Base
// @Failure 503 {object} response.Base
// @Router /health [get]
func (h *HTTP) HealthCheck(c *fiber.Ctx) error {
	if err := h.DB.Read.Ping(); err != nil {
		logger.ErrorWithStack(err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "SERVER UNHEALTHY"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}
