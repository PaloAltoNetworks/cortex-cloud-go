// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/stretchr/testify/require"
)

const (
	TestAPIKey   = "test-key"
	TestAPIKeyID = 123
)

func setupTest(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	config := &client.Config{
		ApiUrl:    server.URL,
		ApiKey:    "test-key",
		ApiKeyId:  123,
		Transport: server.Client().Transport.(*http.Transport),
	}
	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.internalClient)
	return client, server
}
