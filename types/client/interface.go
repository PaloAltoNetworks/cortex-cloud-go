package types

import (
	"context"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
)

type CortexClient interface {
	IsCortexClient()
	ValidateAPIKey(ctx context.Context) (bool, error)
	FQDN() string
	APIURL() string
	APIKeyType() string
	SkipSSLVerify() bool
	Timeout() time.Duration
	MaxRetries() int
	RetryMaxDelay() time.Duration
	CrashStackDir() string
	LogLevel() string
	Logger() log.Logger
	SkipLoggingTransport() bool
}
