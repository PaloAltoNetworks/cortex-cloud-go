// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
)

// API endpoint path specification.
const (
	// Cloud Account Management
	ListAccountsByInstanceEndpoint           = "public_api/v1/cloud_onboarding/get_accounts"
	EnableDisableAccountsInInstancesEndpoint = "public_api/v1/cloud_onboarding/enable_disable_account"
	// Cloud Integration Instance Management
	CreateIntegrationTemplateEndpoint           = "public_api/v1/cloud_onboarding/create_instance_template"
	GetIntegrationInstanceDetailsEndpoint       = "public_api/v1/cloud_onboarding/get_instance_details"
	ListIntegrationInstancesEndpoint            = "public_api/v1/cloud_onboarding/get_instances"
	EditIntegrationInstanceEndpoint             = "public_api/v1/cloud_onboarding/edit_instance"
	EnableOrDisableIntegrationInstancesEndpoint = "public_api/v1/cloud_onboarding/enable_disable_instance"
	DeleteIntegrationInstancesEndpoint          = "public_api/v1/cloud_onboarding/delete_instance"
	// General
	GetAzureApprovedTenantsEndpoint = "public_api/v1/cloud_onboarding/get_azure_approved_tenants"
	// Outpost Management
	CreateOutpostTemplateEndpoint = "public_api/v1/cloud_onboarding/create_outpost_template"
	UpdateOutpostEndpoint         = "public_api/v1/cloud_onboarding/edit_outpost"
	ListOutpostsEndpoint          = "public_api/v1/cloud_onboarding/get_outposts"
)

// Option is a functional option for configuring the client.
type Option = config.Option

var (
	// WithCortexFQDN is an option to set the Cortex FQDN.
	WithCortexFQDN = config.WithCortexFQDN
	// WithCortexAPIURL is an option to set the Cortex API URL.
	WithCortexAPIURL = config.WithCortexAPIURL
	// WithCortexAPIKey is an option to set the Cortex API key.
	WithCortexAPIKey = config.WithCortexAPIKey
	// WithCortexAPIKeyID is an option to set the Cortex API key ID.
	WithCortexAPIKeyID = config.WithCortexAPIKeyID
	// WithCortexAPIKeyType is an option to set the Cortex API key type.
	WithCortexAPIKeyType = config.WithCortexAPIKeyType
	// WithCortexAPIPort is an option to set the Cortex API port.
	WithHeaders = config.WithHeaders
	// WithAgent is an option to set the user agent.
	WithAgent = config.WithAgent
	// WithSkipSSLVerify is an option to skip TLS certificate verification.
	WithSkipSSLVerify = config.WithSkipSSLVerify
	// WithTransport is an option to set the HTTP transport.
	WithTransport = config.WithTransport
	// WithTimeout is an option to set the HTTP timeout.
	WithTimeout = config.WithTimeout
	// WithMaxRetries is an option to set the maximum number of retries.
	WithMaxRetries = config.WithMaxRetries
	// WithRetryMaxDelay is an option to set the maximum retry delay.
	WithRetryMaxDelay = config.WithRetryMaxDelay
	// WithCrashStackDir is an option to set the crash stack directory.
	WithCrashStackDir = config.WithCrashStackDir
	// WithLogLevel is an option to set the log level.
	WithLogLevel = config.WithLogLevel
	// WithLogger is an option to set the logger.
	WithLogger = config.WithLogger
	// WithSkipLoggingTransport is an option to skip logging transport.
	WithSkipLoggingTransport = config.WithSkipLoggingTransport
)

// Client is the client for the namespace.
type Client struct {
	internalClient *client.Client
}

// Marker method for CortexClient interface compliance.
func (Client) IsCortexClient() {}

// ValidateAPIKey validates the configured API Key against the target
// Cortex tenant.
func (c *Client) ValidateAPIKey(ctx context.Context) (bool, error) {
	return c.internalClient.ValidateAPIKey(ctx)
}

// NewClient returns a new client for this namespace.
func NewClient(opts ...Option) (*Client, error) {
	cfg := config.NewConfig(opts...)
	internalClient, err := client.NewClientFromConfig(cfg)
	return &Client{internalClient: internalClient}, err
}

// NewClientFromFile creates a new client from a configuration file.
func NewClientFromFile(filepath string, checkEnvironment bool) (*Client, error) {
	config, err := config.NewConfigFromFile(filepath, checkEnvironment)
	if err != nil {
		return nil, err
	}
	return NewClient(config.GetOptions()...)
}

// FQDN returns the FQDN of the Cortex tenant.
func (c *Client) FQDN() string { return c.internalClient.FQDN() }

// APIURL returns the API URL for the Cortex.
func (c *Client) APIURL() string { return c.internalClient.APIURL() }

// APIKeyType returns the Cortex API key type.
func (c *Client) APIKeyType() string { return c.internalClient.APIKeyType() }

// SkipSSLVerify returns whether to skip TLS certificate verification.
func (c *Client) SkipSSLVerify() bool { return c.internalClient.SkipSSLVerify() }

// Timeout returns the HTTP timeout.
func (c *Client) Timeout() time.Duration { return c.internalClient.Timeout() }

// MaxRetries returns the maximum number of retries.
func (c *Client) MaxRetries() int { return c.internalClient.MaxRetries() }

// RetryMaxDelay returns the maximum retry delay.
func (c *Client) RetryMaxDelay() time.Duration { return c.internalClient.RetryMaxDelay() }

// CrashStackDir returns the crash stack directory.
func (c *Client) CrashStackDir() string { return c.internalClient.CrashStackDir() }

// LogLevel returns the log level.
func (c *Client) LogLevel() string { return c.internalClient.LogLevel() }

// Logger returns the logger.
func (c *Client) Logger() log.Logger { return c.internalClient.Logger() }

// SkipLoggingTransport returns whether to skip logging transport.
func (c *Client) SkipLoggingTransport() bool { return c.internalClient.SkipLoggingTransport() }
