// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("should create new client with valid config", func(t *testing.T) {
		client, err := NewClient(
			WithCortexAPIURL("https://api.example.com"),
			WithCortexAPIKey("test-key"),
			WithCortexAPIKeyID(123),
		)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.internalClient)
	})
}

func TestNewClientFromFile(t *testing.T) {
	t.Run("should create new client from file", func(t *testing.T) {
		// Create a temporary config file
		content := []byte(`{
			"api_url": "https://api.from.file",
			"api_key": "key-from-file",
			"api_key_id": 456
		}`)
		tmpfile, err := os.CreateTemp("", "test-config-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		// Create client from file
		client, err := NewClientFromFile(tmpfile.Name(), false)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.internalClient)
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		client, err := NewClientFromFile("/non/existent/file.json", false)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
