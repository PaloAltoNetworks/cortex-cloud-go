package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfo(t *testing.T) {
	t.Skip("Skipping build info test for now.")
}

func TestNewClient(t *testing.T) {
	t.Run("should return error for nil config from NewClientFromConfig", func(t *testing.T) {
		client, err := NewClientFromConfig(nil)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, "received nil Config", err.Error())
	})

	t.Run("should create new client with valid config from NewClientFromConfig", func(t *testing.T) {
		cfg := config.NewConfig(
			config.WithCortexAPIURL("https://api.example.com"),
			config.WithCortexAPIKey("test-key"),
			config.WithCortexAPIKeyID(123),
		)
		client, err := NewClientFromConfig(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, cfg, client.config)
	})

	t.Run("should use default logger if none is provided", func(t *testing.T) {
		cfg := config.NewConfig(
			config.WithCortexAPIURL("https://api.example.com"),
			config.WithCortexAPIKey("test-key"),
			config.WithCortexAPIKeyID(123),
		)
		client, err := NewClientFromConfig(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config.Logger())
	})
}

func TestGenerateHeaders(t *testing.T) {
	cfg := config.NewConfig(
		config.WithCortexAPIURL("https://api.example.com"),
		config.WithCortexAPIKey("test-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCheckEnvironment(false),
		config.WithAgent("test-agent"),
	)
	client, _ := NewClientFromConfig(cfg)

	t.Run("should generate headers with content type", func(t *testing.T) {
		headers, err := client.generateHeaders(true)
		assert.NoError(t, err)
		assert.Equal(t, "application/json", headers["Content-Type"])
		assert.Equal(t, "test-agent", headers["User-Agent"])
		assert.Equal(t, "1", headers["x-xdr-auth-id"])
		assert.NotEmpty(t, headers["x-xdr-nonce"])
		assert.NotEmpty(t, headers["x-xdr-timestamp"])
		assert.NotEmpty(t, headers["Authorization"])
	})

	t.Run("should generate headers without content type", func(t *testing.T) {
		headers, err := client.generateHeaders(false)
		assert.NoError(t, err)
		assert.NotContains(t, headers, "Content-Type")
		assert.Equal(t, "test-agent", headers["User-Agent"])
	})
}

func TestBuildRequestURL(t *testing.T) {
	cfg := config.NewConfig(
		config.WithCortexAPIURL("https://server.com/api/"),
		config.WithCortexAPIKey("key"),
		config.WithCortexAPIKeyID(1),
		config.WithCheckEnvironment(false),
	)
	client, _ := NewClientFromConfig(cfg)

	t.Run("should build url with path and query params", func(t *testing.T) {
		endpoint := "v1/resource"
		pathParams := &[]string{"12345"}
		queryParams := &url.Values{"key": []string{"value"}}
		expectedURL := "https://server.com/api/v1/resource/12345?key=value"

		actualURL, err := client.buildRequestURL(endpoint, pathParams, queryParams)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL, actualURL)
	})

	t.Run("should handle no params", func(t *testing.T) {
		endpoint := "v1/health"
		expectedURL := "https://server.com/api/v1/health"

		actualURL, err := client.buildRequestURL(endpoint, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedURL, actualURL)
	})

	t.Run("should return error for invalid base url", func(t *testing.T) {
		badCfg := config.NewConfig(config.WithCortexAPIURL("::not-a-url"))
		badClient, _ := NewClientFromConfig(badCfg)
		_, err := badClient.buildRequestURL("v1/endpoint", nil, nil)
		assert.Error(t, err)
	})
}

func TestIsRetryableHTTPStatus(t *testing.T) {
	assert.True(t, isRetryableHTTPStatus(http.StatusUnauthorized))
	assert.True(t, isRetryableHTTPStatus(http.StatusTooManyRequests))
	assert.True(t, isRetryableHTTPStatus(http.StatusBadGateway))
	assert.True(t, isRetryableHTTPStatus(http.StatusServiceUnavailable))
	assert.True(t, isRetryableHTTPStatus(http.StatusGatewayTimeout))
	assert.False(t, isRetryableHTTPStatus(http.StatusOK))
	assert.False(t, isRetryableHTTPStatus(http.StatusBadRequest))
	assert.False(t, isRetryableHTTPStatus(http.StatusNotFound))
	assert.False(t, isRetryableHTTPStatus(http.StatusInternalServerError))
}

func TestDo(t *testing.T) {
	cfg := config.NewConfig(
		config.WithCortexAPIURL("https://testing.com"),
		config.WithCortexAPIKey("key"),
		config.WithCortexAPIKeyID(1),
		config.WithMaxRetries(1),
	)

	t.Run("should succeed on first try", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status":"success"}`)),
		}
		client.testData = []*http.Response{mockResponse}

		var output map[string]string
		_, err := client.Do(context.Background(), "POST", "test", nil, nil, nil, &output, nil)

		assert.NoError(t, err)
		assert.Equal(t, "success", output["status"])
		assert.Equal(t, 1, client.testIndex)
	})

	t.Run("should retry on retryable error and then succeed", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		retryResponse := &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		successResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status":"ok"}`)),
		}
		client.testData = []*http.Response{retryResponse, successResponse}

		var output map[string]string
		_, err := client.Do(context.Background(), "GET", "test", nil, nil, nil, &output, nil)

		assert.NoError(t, err)
		assert.Equal(t, "ok", output["status"])
		assert.Equal(t, 2, client.testIndex)
	})

	t.Run("should fail on non-retryable error", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		errorResponse := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{"err_code":404,"err_msg":"Not Found"}`)),
		}
		client.testData = []*http.Response{errorResponse}

		_, err := client.Do(context.Background(), "GET", "test", nil, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, 1, client.testIndex)
	})

	t.Run("should fail after max retries", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		retryResponse := &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		// Provide more error responses than max retries
		client.testData = []*http.Response{retryResponse, retryResponse, retryResponse}

		_, err := client.Do(context.Background(), "GET", "test", nil, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, 2, client.testIndex) // 1 initial + 1 retry
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		_, err := client.Do(ctx, "GET", "test", nil, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request cancelled by context")
	})

	t.Run("should wrap request and unwrap response", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		// Mock response needs to be a ReadCloser
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data":{"status":"success"}}`)),
		}
		client.testData = []*http.Response{mockResponse}
		client.testIndex = 0 // Reset test index

		var output map[string]string
		input := map[string]string{"hello": "world"}
		opts := &DoOptions{
			RequestWrapperKeys:  []string{"request"},
			ResponseWrapperKeys: []string{"data"},
		}
		_, err := client.Do(context.Background(), "POST", "test", nil, nil, &input, &output, opts)

		assert.NoError(t, err)
		assert.Equal(t, "success", output["status"])
		assert.Equal(t, 1, client.testIndex)
	})

	t.Run("should wrap request and unwrap response with multiple keys", func(t *testing.T) {
		client, _ := NewClientFromConfig(cfg)
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data":{"payload":{"status":"success"}}}`)),
		}
		client.testData = []*http.Response{mockResponse}
		client.testIndex = 0 // Reset test index

		var output map[string]string
		input := map[string]string{"hello": "world"}
		opts := &DoOptions{
			RequestWrapperKeys:  []string{"request", "payload"},
			ResponseWrapperKeys: []string{"data", "payload"},
		}
		_, err := client.Do(context.Background(), "POST", "test", nil, nil, &input, &output, opts)

		assert.NoError(t, err)
		assert.Equal(t, "success", output["status"])
		assert.Equal(t, 1, client.testIndex)
	})
}