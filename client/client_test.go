// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	//"runtime"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfo(t *testing.T) {
	// This test is a placeholder and will likely fail without build-time variables.
	// It can be adjusted based on the actual build process.
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
		config := &Config{
			ApiUrl:   "https://api.example.com",
			ApiKey:   "test-key",
			ApiKeyId: 123,
		}
		client, err := NewClientFromConfig(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, config, client.config)
	})

	t.Run("should create new client with valid values from NewClient", func(t *testing.T) {
		client, err := NewClient("https://api.example.com", "test-key", 123, false)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, "https://api.example.com", client.config.ApiUrl)
	})

	t.Run("should use default logger if none is provided", func(t *testing.T) {
		config := &Config{
			ApiUrl:   "https://api.example.com",
			ApiKey:   "test-key",
			ApiKeyId: 123,
		}
		client, err := NewClientFromConfig(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config.Logger)
	})
}

func TestGenerateHeaders(t *testing.T) {
	client, _ := NewClient("https://api.example.com", "test-api-key", 1, false, WithAgent("test-agent"))

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
	client, _ := NewClient("https://server.com/api/", "key", 1, false)

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
		client.config.ApiUrl = "::not-a-url"
		_, err := client.buildRequestURL("v1/endpoint", nil, nil)
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
	config := &Config{
		ApiUrl:     "https://testing.com",
		ApiKey:     "key",
		ApiKeyId:   1,
		MaxRetries: 1,
	}

	t.Run("should succeed on first try", func(t *testing.T) {
		client, _ := NewClientFromConfig(config)
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
		client, _ := NewClientFromConfig(config)
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
		client, _ := NewClientFromConfig(config)
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
		client, _ := NewClientFromConfig(config)
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
		client, _ := NewClientFromConfig(config)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		_, err := client.Do(ctx, "GET", "test", nil, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request cancelled by context")
	})

	t.Run("should wrap request and unwrap response", func(t *testing.T) {
		client, _ := NewClientFromConfig(config)
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
		client, _ := NewClientFromConfig(config)
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
