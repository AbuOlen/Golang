package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestService_AddUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful user creation", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		// Mock successful insert
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		user, err := svc.AddUser(context.Background(), "testuser", "hashedpassword")
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "testuser", user.Name)
		assert.Equal(t, "hashedpassword", user.PwdHash)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.CreatedAt.IsZero())
	})

	mt.Run("duplicate user", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		// Mock duplicate key error
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		user, err := svc.AddUser(context.Background(), "existinguser", "hashedpassword")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
		assert.Nil(t, user)
	})
}

func TestService_GetUserByName(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("user found", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		expectedUser := &User{
			ID:        "testid",
			Name:      "existinguser",
			PwdHash:   "hashedpassword",
			CreatedAt: time.Now().UTC(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"id", expectedUser.ID},
			{"name", expectedUser.Name},
			{"pwd_hash", expectedUser.PwdHash},
			{"created_at", expectedUser.CreatedAt},
		}))

		user, err := svc.GetUserByName(context.Background(), expectedUser.Name)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.PwdHash, user.PwdHash)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Message: "not found",
			}))

		user, err := svc.GetUserByName(context.Background(), "nonexistentuser")
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestService_GetUserByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("user found", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		expectedUser := &User{
			ID:        "testid123",
			Name:      "testuser",
			PwdHash:   "hashedpassword",
			CreatedAt: time.Now().UTC(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"id", expectedUser.ID},
			{"name", expectedUser.Name},
			{"pwd_hash", expectedUser.PwdHash},
			{"created_at", expectedUser.CreatedAt},
		}))

		user, err := svc.GetUserByID(context.Background(), expectedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.PwdHash, user.PwdHash)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		svc := &Service{
			userColl: mt.Coll,
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Message: "not found",
			}))

		user, err := svc.GetUserByID(context.Background(), "nonexistentid")
		require.Error(t, err)
		assert.Nil(t, user)
	})
}
