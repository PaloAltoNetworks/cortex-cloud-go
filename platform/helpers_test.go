package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.NotNil(t, client)
	return &Client{
		internalClient: client, 
	}, server
}
