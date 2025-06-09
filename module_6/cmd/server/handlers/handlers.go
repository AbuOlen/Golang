package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"module6/cmd/server/handlers/auth"
	"module6/cmd/server/handlers/counter"
	"module6/cmd/server/handlers/health"
	"module6/cmd/server/handlers/message"
	"module6/cmd/server/middlewares"
	"module6/internal/config"
	"module6/internal/services"
)

type Handlers struct {
	Health  *health.Handler
	Counter *counter.Handler
	Auth    *auth.Handler
	Message *message.Handler

	mdlwr *middlewares.Middlewares
}

func NewHandlers(cfg *config.Config, svcs *services.Services, mdlwr *middlewares.Middlewares) *Handlers {
	return &Handlers{
		Health:  health.NewHandler(cfg),
		Counter: counter.NewHandler(svcs),
		Auth:    auth.NewHandler(svcs),
		Message: message.NewHandler(svcs, cfg),
		mdlwr:   mdlwr,
	}
}

func (h *Handlers) RegisterRoutes(router fiber.Router) {

	router.Get("/metrics", prometheusHandler())
	router.Get("/health", h.Health.Health)
	router.Get("/docs/*", swagger.HandlerDefault)

	api := router.Group("/api")
	api.Use(h.mdlwr.Telemetry.Handle)
	api.Use(h.mdlwr.Logger.Handle)

	authGroup := api.Group("/auth")
	authGroup.Post("/signup", h.Auth.SignUp)
	authGroup.Post("/signin", h.Auth.SignIn)

	messageGroup := api.Group("/message")
	messageGroup.Use(h.mdlwr.Auth.Handle)
	messageGroup.Post("/send", h.Message.SendMessage)
	messageGroup.Post("/list", h.Message.List)
	messageGroup.Get("/listen/:channel_name", websocket.New(h.Message.Listen))
}

func prometheusHandler() fiber.Handler {
	handler := promhttp.Handler()
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(handler)(c.Context())
		return nil
	}
}
