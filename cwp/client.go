// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
)

const (
	ListPoliciesEndpoint  = "public_api/v1/cwp/get_policies"
	CreatePolicyEndpoint  = "public_api/v1/cwp/policies"
	GetPolicyByIDEndpoint = "public_api/v1/cwp/get_policy_details"
	EditPolicyEndpoint    = "public_api/v1/cwp/edit_policy"
	DeletePolicyEndpoint  = "public_api/v1/cwp/delete_policy"
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
