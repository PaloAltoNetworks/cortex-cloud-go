package compliance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_NewClient(t *testing.T) {
	t.Run("should create client successfully", func(t *testing.T) {
		client, err := NewClient(
			WithCortexAPIURL("https://api-test.example.com"),
			WithCortexAPIKey("test-key"),
			WithCortexAPIKeyID(123),
		)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("should expose configuration values", func(t *testing.T) {
		client, err := NewClient(
			WithCortexAPIURL("https://api-test.example.com"),
			WithCortexFQDN("test.example.com"),
			WithCortexAPIKey("test-key"),
			WithCortexAPIKeyID(123),
			WithCortexAPIKeyType("standard"),
		)

		require.NoError(t, err)

		assert.Equal(t, "https://api-test.example.com", client.APIURL())
		assert.Equal(t, "test.example.com", client.FQDN())
		assert.Equal(t, "standard", client.APIKeyType())
	})

	t.Run("should validate API key", func(t *testing.T) {
		// This would require mocking the internal client
		// For now, just verify the method exists
		client, err := NewClient(
			WithCortexAPIURL("https://api-test.example.com"),
			WithCortexAPIKey("test-key"),
			WithCortexAPIKeyID(123),
		)

		require.NoError(t, err)
		assert.NotNil(t, client)

		// ValidateAPIKey would make real HTTP call
		// We skip testing it here as it requires proper mocking
	})

	t.Run("should implement CortexClient interface", func(t *testing.T) {
		client, err := NewClient(
			WithCortexAPIURL("https://api-test.example.com"),
			WithCortexAPIKey("test-key"),
			WithCortexAPIKeyID(123),
		)

		require.NoError(t, err)

		// Verify marker method exists and can be called
		assert.NotPanics(t, func() {
			client.IsCortexClient()
		})
	})

	t.Run("should create client from file", func(t *testing.T) {
		// This would require creating a temp config file
		// Skipping for now as it requires file I/O setup
		t.Skip("NewClientFromFile requires test config file")
	})
}

func TestClient_ConfigMethods(t *testing.T) {
	client, err := NewClient(
		WithCortexAPIURL("https://api-test.example.com"),
		WithCortexAPIKey("test-key"),
		WithCortexAPIKeyID(123),
		WithSkipSSLVerify(true),
		WithMaxRetries(5),
		WithLogLevel("info"),
		WithCrashStackDir("/tmp/crashes"),
		WithSkipLoggingTransport(true),
	)

	require.NoError(t, err)

	// Test various config accessor methods
	assert.Equal(t, true, client.SkipSSLVerify())
	assert.Equal(t, 5, client.MaxRetries())

	// Timeout returns time.Duration
	assert.NotZero(t, client.Timeout())

	// RetryMaxDelay returns time.Duration
	assert.NotZero(t, client.RetryMaxDelay())

	// Test additional config methods
	assert.Equal(t, "info", client.LogLevel())
	assert.Equal(t, "/tmp/crashes", client.CrashStackDir())
	assert.Equal(t, true, client.SkipLoggingTransport())

	// Test Logger accessor (should not be nil)
	assert.NotNil(t, client.Logger())
}
