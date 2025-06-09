package message

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"module6/cmd/server/handlers/message/telemetry"
	"module6/cmd/server/middlewares/auth"
	"module6/cmd/server/parser"
	"module6/internal/config"
	"module6/internal/logger"
	"module6/internal/services"
	ms "module6/internal/services/message"
	"time"
)

type Handler struct {
	svc *services.Services
	cfg *config.Config
}

func NewHandler(svc *services.Services, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

type SendMessageRequestBody struct {
	Channel string `json:"channel" validate:"required,min=5,max=20"`
	Message string `json:"message" validate:"required,min=1,max=1024"`
}

type SendMessageResponseBody struct {
	Message *ms.Message `json:"message"`
}

// @Summary SendMessage
// @Description Sending new message to channel
// @Tags message
// @Param body body SendMessageRequestBody true "Send message request body"
// @Success 200 {object} SendMessageResponseBody "Send message succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 500 {object} any "Internal server error"
// @Router /api/message/send [post]
func (h *Handler) SendMessage(ctx *fiber.Ctx) error {
	msg := &SendMessageRequestBody{}
	err := parser.ParseBody(ctx, msg)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u := auth.MustGetUser(ctx)

	m, err := h.svc.Message.Send(ctx.Context(), u.ID, msg.Channel, msg.Message)
	if err != nil {
		telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusFailure)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	telemetry.AddMessageSend(h.cfg.InstanceName, telemetry.MessageStatusSuccess)

	return ctx.JSON(&SendMessageResponseBody{
		Message: m,
	})
}

type ListMessagesRequestBody struct {
	Channel        string  `json:"channel" validate:"required,min=5,max=20"`
	LastReceivedID *string `json:"last_received_id"`
}

type ListMessagesResponseBody struct {
	Messages []*ms.Message `json:"messages"`
}

// List swagger
// @Summary List
// @Description List of messages by channel
// @Tags message
// @Param body body ListMessagesRequestBody true "List messages request body"
// @Success 200 {object} ListMessagesResponseBody "List messages retrieval succeeded"
// @Failure 400 {object} any "Client error"
// @Failure 500 {object} any "Internal server error"
// @Router /api/message/list [post]
func (h *Handler) List(ctx *fiber.Ctx) error {
	listParams := &ListMessagesRequestBody{}
	err := parser.ParseBody(ctx, listParams)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	m, err := h.svc.Message.GetList(ctx.Context(), ms.ListParams{
		Channel:        listParams.Channel,
		LastReceivedID: listParams.LastReceivedID,
		Limit:          2,
	})
	if err != nil {
		telemetry.AddMessageList(h.cfg.InstanceName, telemetry.MessageStatusFailure)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	telemetry.AddMessageList(h.cfg.InstanceName, telemetry.MessageStatusSuccess)
	return ctx.JSON(m)
}

func (h *Handler) Listen(conn *websocket.Conn) {
	channelName := conn.Params("channel_name")
	u := auth.MustGetUserWs(conn)
	logger.GetLogger().Info(context.Background(), "user start listening", slog.String("user_id", u.ID), slog.String("channel_name", channelName))

	// Setup pong handler and read deadline
	err := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	})
	cc := make(chan struct{})
	// Send pings periodically
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			logger.GetLogger().Info(context.Background(), "ticker tick")
			if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second)); err != nil {
				_ = conn.Close()
				logger.GetLogger().Info(context.Background(), "no response from websocket")
				cc <- struct{}{}
				return
			}
		}
	}()

	ch := make(chan *ms.Message)
	h.svc.Message.SubscribeToChannel(channelName, ch)
	defer h.svc.Message.UnsubscribeFromChannel(channelName, ch)

	for {
		select {
		case m := <-ch:
			err := conn.WriteJSON(m)
			if err != nil {
				_ = conn.Close()
				return
			}
		case <-cc:
			return
		}
	}

}
