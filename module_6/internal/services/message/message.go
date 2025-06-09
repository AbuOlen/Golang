package message

import "time"

// Message send in chats
// swagger:model
type Message struct {
	ID        string    `json:"id" bson:"id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Channel   string    `json:"channel" bson:"channel"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type ListParams struct {
	Channel        string  `json:"channel"`
	Limit          int64   `json:"limit"`
	LastReceivedID *string `json:"last_received_id" bson:"last_received_id"`
}
