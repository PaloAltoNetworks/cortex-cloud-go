// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
	"github.com/stretchr/testify/require"
)

const (
	TestAPIKey   = "test-key"
	TestAPIKeyID = 123
)

func setupTest(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	t.Logf("running setupTest for %s", t.Name())
	server := httptest.NewServer(handler)
	config := &client.Config{
		ApiUrl:    server.URL,
		ApiKey:    "test-key",
		ApiKeyId:  123,
		Transport: server.Client().Transport.(*http.Transport),
		Logger: log.TflogAdapter{},
	}
	client, err := client.NewClientFromConfig(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	return &Client{
		internalClient: client, 
	}, server
}
