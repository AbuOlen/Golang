package message

import (
	"context"
	"fmt"
	"github.com/oklog/ulid/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"module6/internal/clients"
	"module6/internal/constants"
	"time"
)

type Service struct {
	messagesColl *mongo.Collection
	ecs          cmap.ConcurrentMap[string, *EventChannel]
}

func NewService(c *clients.Clients) *Service {
	return &Service{
		messagesColl: c.Mongo.Db.Collection(constants.CollectionMessages),
		ecs:          cmap.New[*EventChannel](),
	}
}

func (s *Service) Send(ctx context.Context, userID, channel, message string) (*Message, error) {
	utcNow := time.Now().UTC()

	m := &Message{
		ID:        ulid.MustNew(uint64(utcNow.Unix()), ulid.DefaultEntropy()).String(),
		UserID:    userID,
		Channel:   channel,
		Message:   message,
		CreatedAt: utcNow,
	}

	_, err := s.messagesColl.InsertOne(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("failed to send message %w", err)
	}
	ec, ok := s.ecs.Get(m.Channel)
	if ok {
		ec.PublishMessage(m)
	}

	return m, nil
}

func (s *Service) GetList(ctx context.Context, listParams ListParams) ([]Message, error) {
	filter := bson.M{
		"channel": listParams.Channel,
	}
	if listParams.LastReceivedID != nil {
		filter["id"] = bson.M{"$lt": *listParams.LastReceivedID}
	}
	opts := options.Find().SetLimit(listParams.Limit)

	cursor, err := s.messagesColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages %w", err)
	}
	var messages []Message
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from cursor %w", err)
	}
	return messages, nil
}

func (s *Service) SubscribeToChannel(channelName string, c chan<- *Message) {
	ec, ok := s.ecs.Get(channelName)
	if !ok {
		ec = StartEventChannel()
		s.ecs.Set(channelName, ec)
	}
	ec.Subscribe(c)
}

func (s *Service) UnsubscribeFromChannel(channelName string, c chan<- *Message) {
	ec, ok := s.ecs.Get(channelName)
	if !ok {
		return
	}
	ec.Unsubscribe(c)
	if !ec.IsEmpty() {
		return
	}
	ec.Stop()
	s.ecs.Remove(channelName)
}
