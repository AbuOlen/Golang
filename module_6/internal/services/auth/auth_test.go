package auth

import (
	"module6/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestService() *Service {
	cfg := &config.Config{
		TokenTTLMinutes: "60",
	}
	return NewService(cfg)
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Config
		expectedPanic bool
	}{
		{
			name: "valid configuration",
			config: &config.Config{
				TokenTTLMinutes: "60",
			},
			expectedPanic: false,
		},
		{
			name: "invalid TTL minutes",
			config: &config.Config{
				TokenTTLMinutes: "invalid",
			},
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedPanic {
				assert.Panics(t, func() {
					NewService(tt.config)
				})
			} else {
				svc := NewService(tt.config)
				assert.NotNil(t, svc)
				assert.Equal(t, 60*time.Minute, svc.tokenTTLMinutes)
			}
		})
	}
}

func TestService_GeneratePasswordHash(t *testing.T) {
	svc := createTestService()

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "valid password",
			password:    "mySecurePassword123",
			expectError: false,
		},
		{
			name:        "empty password",
			password:    "",
			expectError: false,
		},
		{
			name:        "password with special characters",
			password:    "p@ssw0rd!#$%^&*()",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := svc.GeneratePasswordHash(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)

				// Verify that the hash can be compared
				match, err := svc.CompareHashAndPassword(tt.password, hash)
				assert.NoError(t, err)
				assert.True(t, match)
			}
		})
	}
}

func TestService_CompareHashAndPassword(t *testing.T) {
	svc := createTestService()
	correctPassword := "mySecurePassword123"

	// Generate a hash for testing
	hash, err := svc.GeneratePasswordHash(correctPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectMatch bool
		expectError bool
	}{
		{
			name:        "correct password",
			password:    correctPassword,
			hash:        hash,
			expectMatch: true,
			expectError: false,
		},
		{
			name:        "wrong password",
			password:    "wrongPassword123",
			hash:        hash,
			expectMatch: false,
			expectError: false,
		},
		{
			name:        "invalid hash format",
			password:    correctPassword,
			hash:        "invalid_hash_format",
			expectMatch: false,
			expectError: true,
		},
		{
			name:        "empty password with valid hash",
			password:    "",
			hash:        hash,
			expectMatch: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := svc.CompareHashAndPassword(tt.password, tt.hash)

			if tt.expectError {
				assert.Error(t, err)
				assert.False(t, match)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectMatch, match)
			}
		})
	}
}
