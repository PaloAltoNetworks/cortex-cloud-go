// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
)

// API endpoint path specification.
const (
	// System Management Endpoints
	HealthCheckEndpoint    = "public_api/v1/health_check/"
	GetTenantInfoEndpoint  = "public_api/v1/get_tenant_info/"
	ListUsersEndpoint      = "public_api/v1/rbac/get_users/"
	GetUserGroupEndpoint   = "public_api/v1/rbac/get_user_group/"
	ListRolesEndpoint      = "public_api/v1/rbac/get_roles/"
	SetUserRoleEndpoint    = "public_api/v1/rbac/set_user_role/"
	GetRiskScoreEndpoint   = "public_api/v1/get_risk_score/"
	ListRiskyUsersEndpoint = "public_api/v1/risk/get_risky_users/"
	ListRiskyHostsEndpoint = "public_api/v1/risky_hosts/"

	// Asset Group Endpoints
	CreateAssetGroupEndpoint = "public_api/v1/asset-groups/create"
	UpdateAssetGroupEndpoint = "public_api/v1/asset-groups/update/"
	DeleteAssetGroupEndpoint = "public_api/v1/asset-groups/delete/"
	ListAssetGroupsEndpoint  = "public_api/v1/asset-groups"

	// Auth Settings Endpoints
	ListIDPMetadataEndpoint    = "public_api/v1/sso/get_idp_metadata/"
	ListAuthSettingsEndpoint   = "public_api/v1/sso/get_sso_config/"
	CreateAuthSettingsEndpoint = "public_api/v1/sso/set_config/"
	UpdateAuthSettingsEndpoint = "public_api/v1/sso/set_config/"
	DeleteAuthSettingsEndpoint = "public_api/v1/sso/delete_config/"
)

// Option is a functional option for configuring the client.
type Option = config.Option

var (
	// WithCortexAPIURL is an option to set the Cortex API URL.
	WithCortexAPIURL = config.WithCortexAPIURL
	// WithCortexAPIKey is an option to set the Cortex API key.
	WithCortexAPIKey = config.WithCortexAPIKey
	// WithCortexAPIKeyID is an option to set the Cortex API key ID.
	WithCortexAPIKeyID = config.WithCortexAPIKeyID
	// WithCortexAPIPort is an option to set the Cortex API port.
	WithCortexAPIPort = config.WithCortexAPIPort
	// WithHeaders is an option to set the HTTP headers.
	WithHeaders = config.WithHeaders
	// WithAgent is an option to set the user agent.
	WithAgent = config.WithAgent
	// WithSkipVerifyCertificate is an option to skip TLS certificate verification.
	WithSkipVerifyCertificate = config.WithSkipVerifyCertificate
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
