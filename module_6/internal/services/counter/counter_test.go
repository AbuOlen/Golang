package counter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestService_IncrementOne(t *testing.T) {
	// Create a new mock test
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Create expected response for FindOneAndUpdate
		expectedCount := 42
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"count", expectedCount},
			}},
		})

		// Create service with mock collection
		service := &Service{
			counterColl: mt.Coll,
		}

		// Call the method
		count, err := service.IncrementOne(context.Background())
		require.NoError(t, err)
		assert.Equal(t, expectedCount, count)
	})

	mt.Run("mongodb error", func(mt *mtest.T) {
		// Mock MongoDB error
		mt.AddMockResponses(bson.D{
			{"ok", 0},
			{"errmsg", "operation failed"},
			{"code", 11000},
		})

		service := &Service{
			counterColl: mt.Coll,
		}

		count, err := service.IncrementOne(context.Background())
		require.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "failed to increment counter")
	})

	mt.Run("decode error", func(mt *mtest.T) {
		// Create response with invalid data type
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"count", "not_a_number"}, // This will cause decode error as Count should be int
			}},
		})

		service := &Service{
			counterColl: mt.Coll,
		}

		count, err := service.IncrementOne(context.Background())
		require.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "failed to decode counter result")
	})
}
