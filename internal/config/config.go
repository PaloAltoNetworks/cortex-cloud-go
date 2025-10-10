// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	cortexLog "github.com/PaloAltoNetworks/cortex-cloud-go/log"
)

const (
	CORTEXCLOUD_API_URL_ENV_VAR                 = "CORTEXCLOUD_API_URL"
	CORTEXCLOUD_API_PORT_ENV_VAR                = "CORTEXCLOUD_API_PORT"
	CORTEXCLOUD_API_KEY_ENV_VAR                 = "CORTEXCLOUD_API_KEY"
	CORTEXCLOUD_API_KEY_ID_ENV_VAR              = "CORTEXCLOUD_API_KEY_ID"
	CORTEXCLOUD_HEADERS_ENV_VAR                 = "CORTEXCLOUD_HEADERS"
	CORTEXCLOUD_AGENT_ENV_VAR                   = "CORTEXCLOUD_AGENT"
	CORTEXCLOUD_SKIP_VERIFY_CERTIFICATE_ENV_VAR = "CORTEXCLOUD_SKIP_VERIFY_CERTIFICATE"
	CORTEXCLOUD_CONFIG_FILE_ENV_VAR             = "CORTEXCLOUD_CONFIG_FILE"
	CORTEXCLOUD_TIMEOUT_ENV_VAR                 = "CORTEXCLOUD_TIMEOUT"
	CORTEXCLOUD_MAX_RETRIES_ENV_VAR             = "CORTEXCLOUD_MAX_RETRIES"
	CORTEXCLOUD_RETRY_MAX_DELAY_ENV_VAR         = "CORTEXCLOUD_RETRY_MAX_DELAY"
	CORTEXCLOUD_CRASH_STACK_DIR_ENV_VAR         = "CORTEXCLOUD_CRASH_STACK_DIR"
	CORTEXCLOUD_LOG_LEVEL_ENV_VAR               = "CORTEXCLOUD_LOG_LEVEL"
	CORTEXCLOUD_SKIP_LOGGING_TRANSPORT_ENV_VAR  = "CORTEXCLOUD_SKIP_LOGGING_TRANSPORT"
)

type Config struct {
	checkEnvironmentVars bool
	// TODO: change json tags
	cortexAPIURL          string            `json:"api_url"`
	cortexAPIKey          string            `json:"api_key"`
	cortexAPIKeyID        int               `json:"api_key_id"`
	cortexAPIPort         int               `json:"cortex_api_port"`
	headers               map[string]string `json:"headers"`
	agent                 string            `json:"agent"`
	skipVerifyCertificate bool              `json:"skip_verify_certificate"`
	transport             *http.Transport   `json:"-"`
	timeout               int               `json:"timeout"`
	maxRetries            int               `json:"max_retries"`
	retryMaxDelay         int               `json:"retry_max_delay"`
	crashStackDir         string            `json:"crash_stack_dir"`
	logLevel              string            `json:"log_level"`
	logger                cortexLog.Logger  `json:"-"`
	skipLoggingTransport  bool              `json:"skip_logging_transport"`
}

// CortexAPIURL returns the Cortex API URL.
func (c *Config) CortexAPIURL() string { return c.cortexAPIURL }

// CortexAPIKey returns the Cortex API key.
func (c *Config) CortexAPIKey() string { return c.cortexAPIKey }

// CortexAPIKeyID returns the Cortex API key ID.
func (c *Config) CortexAPIKeyID() int { return c.cortexAPIKeyID }

// CortexAPIPort returns the Cortex API port.
func (c *Config) CortexAPIPort() int { return c.cortexAPIPort }

// Headers returns the HTTP headers.
func (c *Config) Headers() map[string]string { return c.headers }

// Agent returns the user agent.
func (c *Config) Agent() string { return c.agent }

// SkipVerifyCertificate returns whether to skip TLS certificate verification.
func (c *Config) SkipVerifyCertificate() bool { return c.skipVerifyCertificate }

// Transport returns the HTTP transport.
func (c *Config) Transport() *http.Transport { return c.transport }

// Timeout returns the HTTP timeout.
func (c *Config) Timeout() int { return c.timeout }

// MaxRetries returns the maximum number of retries.
func (c *Config) MaxRetries() int { return c.maxRetries }

// RetryMaxDelay returns the maximum retry delay.
func (c *Config) RetryMaxDelay() int { return c.retryMaxDelay }

// CrashStackDir returns the crash stack directory.
func (c *Config) CrashStackDir() string { return c.crashStackDir }

// LogLevel returns the log level.
func (c *Config) LogLevel() string { return c.logLevel }

// Logger returns the logger.
func (c *Config) Logger() cortexLog.Logger { return c.logger }

// SkipLoggingTransport returns whether to skip logging transport.
func (c *Config) SkipLoggingTransport() bool { return c.skipLoggingTransport }

func NewConfig(opts ...Option) *Config {
	config := &Config{
		checkEnvironmentVars:  true,
		cortexAPIPort:         443,
		headers:               make(map[string]string),
		agent:                 "",
		skipVerifyCertificate: false,
		transport:             http.DefaultTransport.(*http.Transport),
		timeout:               30, // 30 seconds
		maxRetries:            3,
		retryMaxDelay:         60, // 60 seconds
		crashStackDir:         os.TempDir(),
		logLevel:              "info",
		logger:                nil,
		skipLoggingTransport:  false,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.checkEnvironmentVars {
		config.overwriteFromEnvVars()
	}

	return config
}

func NewConfigFromFile(filepath string, checkEnvironment bool) (*Config, error) {
	cBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration file: %s", err)
	}

	var cFile Config
	if err = json.Unmarshal(cBytes, &cFile); err != nil {
		return nil, fmt.Errorf("Error unmarshalling configuration file: %s", err)
	}

	return NewConfig(
		WithCortexAPIURL(cFile.cortexAPIURL),
		WithCortexAPIKey(cFile.cortexAPIKey),
		WithCortexAPIKeyID(cFile.cortexAPIKeyID),
		WithCheckEnvironment(checkEnvironment),
		WithCortexAPIPort(cFile.cortexAPIPort),
		WithHeaders(cFile.headers),
		WithAgent(cFile.agent),
		WithSkipVerifyCertificate(cFile.skipVerifyCertificate),
		WithTransport(cFile.transport),
		WithTimeout(cFile.timeout),
		WithMaxRetries(cFile.maxRetries),
		WithRetryMaxDelay(cFile.retryMaxDelay),
		WithCrashStackDir(cFile.crashStackDir),
		WithLogLevel(cFile.logLevel),
		WithLogger(cFile.logger),
		WithSkipLoggingTransport(cFile.skipLoggingTransport),
	), nil
}

func (c *Config) GetOptions() []Option {
	return []Option{
		WithCortexAPIURL(c.cortexAPIURL),
		WithCortexAPIKey(c.cortexAPIKey),
		WithCortexAPIKeyID(c.cortexAPIKeyID),
		WithCortexAPIPort(c.cortexAPIPort),
		WithHeaders(c.headers),
		WithAgent(c.agent),
		WithSkipVerifyCertificate(c.skipVerifyCertificate),
		WithTransport(c.transport),
		WithTimeout(c.timeout),
		WithMaxRetries(c.maxRetries),
		WithRetryMaxDelay(c.retryMaxDelay),
		WithCrashStackDir(c.crashStackDir),
		WithLogLevel(c.logLevel),
		WithLogger(c.logger),
		WithSkipLoggingTransport(c.skipLoggingTransport),
	}
}

func (c Config) Validate() error {
	// TODO
	// - Make sure URL begins with `api-`
	// - Hit healthcheck endpoint
	// - Check if basic or advanced API key
	return nil
}

// SetDefaults sets default values for the configuration.
func (c *Config) SetDefaults() {
	if c.logger == nil {
		c.logger = cortexLog.DefaultLogger{Logger: log.Default()}
	}
}

func (c *Config) overwriteFromEnvVars() {
	if envApiUrl, ok := os.LookupEnv(CORTEXCLOUD_API_URL_ENV_VAR); ok {
		c.cortexAPIURL = envApiUrl
	}

	if envApiKey, ok := os.LookupEnv(CORTEXCLOUD_API_KEY_ENV_VAR); ok {
		c.cortexAPIKey = envApiKey
	}

	if envApiKeyId, ok := os.LookupEnv(CORTEXCLOUD_API_KEY_ID_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envApiKeyId); err == nil {
			c.cortexAPIKeyID = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEXCLOUD_API_KEY_ID_ENV_VAR, envApiKeyId)
		}
	}

	if envApiPort, ok := os.LookupEnv(CORTEXCLOUD_API_PORT_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envApiPort); err == nil {
			c.cortexAPIPort = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEXCLOUD_API_PORT_ENV_VAR, envApiPort)
		}
	}

	if envHeaders, ok := os.LookupEnv(CORTEXCLOUD_HEADERS_ENV_VAR); ok {
		// Example: HEADERS="Content-Type=application/json,Authorization=Bearer xyz"
		if c.headers == nil {
			c.headers = make(map[string]string)
		}

		for pair := range strings.SplitSeq(envHeaders, ",") {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				c.headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	if envAgent, ok := os.LookupEnv(CORTEXCLOUD_AGENT_ENV_VAR); ok {
		c.agent = envAgent
	}

	if envSkipVerifyCertificate, ok := os.LookupEnv(CORTEXCLOUD_SKIP_VERIFY_CERTIFICATE_ENV_VAR); ok {
		if parsedBool, err := strconv.ParseBool(envSkipVerifyCertificate); err == nil {
			c.skipVerifyCertificate = parsedBool
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected true/false.\n", CORTEXCLOUD_SKIP_VERIFY_CERTIFICATE_ENV_VAR, envSkipVerifyCertificate)
		}
	}

	if envTimeout, ok := os.LookupEnv(CORTEXCLOUD_TIMEOUT_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envTimeout); err == nil {
			c.timeout = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEXCLOUD_TIMEOUT_ENV_VAR, envTimeout)
		}
	}

	if envMaxRetries, ok := os.LookupEnv(CORTEXCLOUD_MAX_RETRIES_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envMaxRetries); err == nil {
			c.maxRetries = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEXCLOUD_MAX_RETRIES_ENV_VAR, envMaxRetries)
		}
	}

	if envRetryMaxDelay, ok := os.LookupEnv(CORTEXCLOUD_RETRY_MAX_DELAY_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envRetryMaxDelay); err == nil {
			c.retryMaxDelay = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEXCLOUD_RETRY_MAX_DELAY_ENV_VAR, envRetryMaxDelay)
		}
	}

	if envCrashStackDir, ok := os.LookupEnv(CORTEXCLOUD_CRASH_STACK_DIR_ENV_VAR); ok {
		c.crashStackDir = envCrashStackDir
	}

	if envLogLevel, ok := os.LookupEnv(CORTEXCLOUD_LOG_LEVEL_ENV_VAR); ok {
		c.logLevel = envLogLevel
	}

	if envSkipLoggingTransport, ok := os.LookupEnv(CORTEXCLOUD_SKIP_LOGGING_TRANSPORT_ENV_VAR); ok {
		if parsedBool, err := strconv.ParseBool(envSkipLoggingTransport); err == nil {
			c.skipLoggingTransport = parsedBool
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected true/false.\n", CORTEXCLOUD_SKIP_LOGGING_TRANSPORT_ENV_VAR, envSkipLoggingTransport)
		}
	}
}
