// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("should return error for nil config", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Error(t, err)
		assert.NotNil(t, client)
		assert.Nil(t, client.internalClient)
		assert.Equal(t, "received nil Config", err.Error()) // This error message comes from the client package now
	})

	t.Run("should create new client with valid config", func(t *testing.T) {
		config := &client.Config{
			ApiUrl:   "https://api.example.com",
			ApiKey:   "test-key",
			ApiKeyId: 123,
		}
		client, err := NewClient(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.internalClient)
	})
}
