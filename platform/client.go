// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/log"
)

// API endpoint path specification.
const (
	// System Management Endpoints
	HealthCheckEndpoint   = "public_api/v1/health_check/"
	GetTenantInfoEndpoint = "public_api/v1/get_tenant_info/"

	ListUsersEndpoint    = "public_api/v1/rbac/get_users/"
	GetUserGroupEndpoint = "public_api/v1/rbac/get_user_group/"

	ListRolesEndpoint    = "public_api/v1/rbac/get_roles/"
	SetUserRoleEndpoint  = "public_api/v1/rbac/set_user_role/"
	GetRiskScoreEndpoint = "public_api/v1/get_risk_score/"

	ListRiskyUsersEndpoint = "public_api/v1/risk/get_risky_users/"
	ListRiskyHostsEndpoint = "public_api/v1/risky_hosts/"

	UserGroupEndpoint = "platform/iam/v1/user-group"

	IamUsersEndpoint = "platform/iam/v1/user"

	// Asset Group Endpoints
	CreateAssetGroupEndpoint = "public_api/v1/asset-groups/create"
	UpdateAssetGroupEndpoint = "public_api/v1/asset-groups/update/"
	DeleteAssetGroupEndpoint = "public_api/v1/asset-groups/delete/"
	ListAssetGroupsEndpoint  = "public_api/v1/asset-groups"

	// Auth Settings Endpoints
	// Updated on Oct 30 per https://cortex-panw.stoplight.io/docs/cortex-cloud/ endpoint change
	ListIDPMetadataEndpoint    = "public_api/v1/authentication-settings/get/metadata"
	ListAuthSettingsEndpoint   = "public_api/v1/authentication-settings/get/settings"
	CreateAuthSettingsEndpoint = "public_api/v1/authentication-settings/create"
	UpdateAuthSettingsEndpoint = "public_api/v1/authentication-settings/update"
	DeleteAuthSettingsEndpoint = "public_api/v1/authentication-settings/delete"
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
	// WithHeaders is an option to set the HTTP headers.
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
