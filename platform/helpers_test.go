package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	t.Logf("running setupTest for %s", t.Name())
	server := httptest.NewServer(handler)
	client, err := NewClient(
		WithCortexAPIURL(server.URL),
		WithCortexAPIKey("test-key"),
		WithCortexAPIKeyID(123),
		WithTransport(server.Client().Transport.(*http.Transport)),
		WithLogger(log.TflogAdapter{}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	return client, server
}
