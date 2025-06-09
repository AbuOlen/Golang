package services

import (
	"module6/internal/clients"
	"module6/internal/config"
	"module6/internal/services/auth"
	"module6/internal/services/counter"
	"module6/internal/services/message"
	"module6/internal/services/user"
)

type Services struct {
	Counter *counter.Service
	Auth    *auth.Service
	User    *user.Service
	Message *message.Service
}

func NewServices(clients *clients.Clients, cfg *config.Config) *Services {
	return &Services{
		Counter: counter.NewService(clients),
		Auth:    auth.NewService(cfg),
		User:    user.NewService(clients),
		Message: message.NewService(clients),
	}
}
