package middlewares

import (
	"module6/cmd/server/middlewares/auth"
	"module6/cmd/server/middlewares/logger"
	"module6/cmd/server/middlewares/telemetry"
	"module6/internal/config"
	"module6/internal/services"
)

type Middlewares struct {
	Logger    *logger.Middleware
	Auth      *auth.Middleware
	Telemetry *telemetry.Middleware
}

func NewMiddlewares(svc *services.Services, cfg *config.Config) *Middlewares {
	return &Middlewares{
		Logger:    logger.NewMiddleware(svc, cfg),
		Auth:      auth.NewMiddleware(svc, cfg),
		Telemetry: telemetry.NewMiddleware(cfg),
	}
}
