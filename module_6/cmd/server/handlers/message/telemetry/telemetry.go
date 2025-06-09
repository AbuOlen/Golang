package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MessageStatus string
type MessageType string

const (
	MessageStatusSuccess MessageStatus = "success"
	MessageStatusFailure MessageStatus = "failure"

	MessageSend MessageType = "send"
	MessageList MessageType = "list"
)

var (
	messagesCounter *prometheus.CounterVec
)

func init() {
	messagesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messages_sent_total",
		Help: "Total number of messages sent",
	}, []string{"instance", "type", "status"})
}

func AddMessageSend(instance string, status MessageStatus) {
	messagesCounter.With(prometheus.Labels{"instance": instance, "type": string(MessageSend), "status": string(status)}).Inc()
}

func AddMessageList(instance string, status MessageStatus) {
	messagesCounter.With(prometheus.Labels{"instance": instance, "type": string(MessageList), "status": string(status)}).Inc()
}
