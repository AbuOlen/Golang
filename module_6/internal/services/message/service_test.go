package message

import (
	"context"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestService_Send(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful message send", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		// Mock successful insert
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		msg, err := svc.Send(context.Background(), "user123", "general", "Hello, World!")
		require.NoError(t, err)
		require.NotNil(t, msg)
		assert.Equal(t, "user123", msg.UserID)
		assert.Equal(t, "general", msg.Channel)
		assert.Equal(t, "Hello, World!", msg.Message)
		assert.NotEmpty(t, msg.ID)
		assert.False(t, msg.CreatedAt.IsZero())
	})

	mt.Run("database error", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		// Mock database error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11600,
			Message: "interrupted at shutdown",
		}))

		msg, err := svc.Send(context.Background(), "user123", "general", "Hello, World!")
		require.Error(t, err)
		assert.Nil(t, msg)
		assert.Contains(t, err.Error(), "failed to send message")
	})
}

func TestService_GetList(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	sampleMessages := []Message{
		{
			ID:        "msg1",
			UserID:    "user1",
			Channel:   "general",
			Message:   "Hello",
			CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			ID:        "msg2",
			UserID:    "user2",
			Channel:   "general",
			Message:   "World",
			CreatedAt: time.Now(),
		},
	}

	mt.Run("successful messages retrieval", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		// Create mock response for Find operation
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"id", sampleMessages[0].ID},
			{"user_id", sampleMessages[0].UserID},
			{"channel", sampleMessages[0].Channel},
			{"message", sampleMessages[0].Message},
			{"created_at", sampleMessages[0].CreatedAt},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"id", sampleMessages[1].ID},
			{"user_id", sampleMessages[1].UserID},
			{"channel", sampleMessages[1].Channel},
			{"message", sampleMessages[1].Message},
			{"created_at", sampleMessages[1].CreatedAt},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		messages, err := svc.GetList(context.Background(), ListParams{
			Channel: "general",
			Limit:   10,
		})

		require.NoError(t, err)
		require.Len(t, messages, 2)
		assert.Equal(t, sampleMessages[0].ID, messages[0].ID)
		assert.Equal(t, sampleMessages[1].ID, messages[1].ID)
	})

	mt.Run("database error", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		// Mock database error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11600,
			Message: "interrupted at shutdown",
		}))

		messages, err := svc.GetList(context.Background(), ListParams{
			Channel: "general",
			Limit:   10,
		})
		require.Error(t, err)
		assert.Nil(t, messages)
		assert.Contains(t, err.Error(), "failed to find messages")
	})

	mt.Run("with last received ID", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		lastID := "msg1"
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"id", "msg2"},
			{"user_id", "user2"},
			{"channel", "general"},
			{"message", "World"},
			{"created_at", time.Now()},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		messages, err := svc.GetList(context.Background(), ListParams{
			Channel:        "general",
			LastReceivedID: &lastID,
			Limit:          10,
		})

		require.NoError(t, err)
		require.Len(t, messages, 1)
		assert.Equal(t, "msg2", messages[0].ID)
	})
}

func TestService_SubscribeAndUnsubscribe(t *testing.T) {
	svc := &Service{
		ecs: cmap.New[*EventChannel](),
	}

	channel := "test-channel"
	msgChan := make(chan *Message)

	// Test Subscribe
	svc.SubscribeToChannel(channel, msgChan)
	ec, exists := svc.ecs.Get(channel)
	assert.True(t, exists)
	assert.NotNil(t, ec)

	// Test Unsubscribe
	svc.UnsubscribeFromChannel(channel, msgChan)
	ec, exists = svc.ecs.Get(channel)
	assert.False(t, exists)
	assert.Nil(t, ec)
}

func TestService_MessageBroadcast(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("message broadcast to subscribers", func(mt *mtest.T) {
		svc := &Service{
			messagesColl: mt.Coll,
			ecs:          cmap.New[*EventChannel](),
		}

		// Create a channel and subscribe to it
		channel := "test-channel"
		msgChan := make(chan *Message, 1)
		svc.SubscribeToChannel(channel, msgChan)

		// Mock successful message send
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Send a message
		sentMsg, err := svc.Send(context.Background(), "user123", channel, "Test broadcast")
		require.NoError(t, err)

		// Check if the message was received by the subscriber
		select {
		case receivedMsg := <-msgChan:
			assert.Equal(t, sentMsg, receivedMsg)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for message broadcast")
		}

		// Cleanup
		svc.UnsubscribeFromChannel(channel, msgChan)
	})
}
