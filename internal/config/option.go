// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	//"log"
	"maps"
	"net/http"

	sdkLog "github.com/PaloAltoNetworks/cortex-cloud-go/log"
)

type Option func(*Config)

// WithCortexAPIURL returns an Option that sets the CortexAPIURL field.
func WithCortexAPIURL(apiUrl string) Option {
	return func(c *Config) {
		c.cortexAPIURL = apiUrl
	}
}

// WithCortexAPIKey returns an Option that sets the CortexAPIKey field.
func WithCortexAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.cortexAPIKey = apiKey
	}
}

// WithCortexAPIKeyID returns an Option that sets the CortexAPIKeyID field.
func WithCortexAPIKeyID(apiKeyId int) Option {
	return func(c *Config) {
		c.cortexAPIKeyID = apiKeyId
	}
}

// WithCheckEnvironment returns an Option that sets the checkEnvironmentVars field.
func WithCheckEnvironment(check bool) Option {
	return func(c *Config) {
		c.checkEnvironmentVars = check
	}
}

// WithCortexAPIPort returns an Option that sets the CortexAPIPort field.
func WithCortexAPIPort(port int) Option {
	return func(c *Config) {
		c.cortexAPIPort = port
	}
}

// WithHeaders returns an Option that sets or adds to the Headers map.
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if c.headers == nil {
			c.headers = make(map[string]string)
		}
		maps.Copy(c.headers, headers)
	}
}

// WithAgent returns an Option that sets the Agent field.
func WithAgent(agent string) Option {
	return func(c *Config) {
		c.agent = agent
	}
}

// WithSkipVerifyCertificate returns an Option that sets the SkipVerifyCertificate field.
func WithSkipVerifyCertificate(skip bool) Option {
	return func(c *Config) {
		c.skipVerifyCertificate = skip
	}
}

// WithTransport returns an Option that sets the Transport field.
func WithTransport(transport *http.Transport) Option {
	return func(c *Config) {
		c.transport = transport
	}
}

// WithTimeout returns an Option that sets the Timeout field (in seconds).
func WithTimeout(timeout int) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithMaxRetries returns an Option that sets the MaxRetries field.
func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.maxRetries = retries
	}
}

// WithRetryMaxDelay returns an Option that sets the RetryMaxDelay field (in seconds).
func WithRetryMaxDelay(delay int) Option {
	return func(c *Config) {
		c.retryMaxDelay = delay
	}
}

// WithCrashStackDir returns an Option that sets the CrashStackDir field.
func WithCrashStackDir(dir string) Option {
	return func(c *Config) {
		c.crashStackDir = dir
	}
}

// WithLogLevel returns an Option that sets the LogLevel field.
func WithLogLevel(level string) Option {
	return func(c *Config) {
		c.logLevel = level
	}
}

// WithLogger returns an Option that sets the Logger field.
func WithLogger(l sdkLog.Logger) Option {
	return func(c *Config) {
		c.logger = l
	}
}

// WithSkipLoggingTransport returns an Option that sets the SkipLoggingTransport field.
func WithSkipLoggingTransport(skip bool) Option {
	return func(c *Config) {
		c.skipLoggingTransport = skip
	}
}
