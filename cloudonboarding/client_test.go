// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"testing"
	"runtime"

	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfo(t *testing.T) {
	expectedGitCommit := "test123"
	expectedGoVersion := runtime.Version()
	expectedBuildDate := "0000-00-00T00:00:00+0000"

	t.Run("should return expected build info", func(t *testing.T) {
		assert.Equal(t, expectedGitCommit, GitCommit)
		assert.Equal(t, expectedGoVersion, GoVersion)
		assert.Equal(t, expectedBuildDate, BuildDate)
	})
}

func TestNewClient(t *testing.T) {
	t.Run("should return error for nil config", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Error(t, err)
		assert.NotNil(t, client)
		assert.Nil(t, client.internalClient)
		assert.Equal(t, "received nil api.Config", err.Error())
	})

	t.Run("should create new client with valid config", func(t *testing.T) {
		config := &api.Config{
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
