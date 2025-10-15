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
	CORTEX_FQDN_ENV_VAR                   = "CORTEX_FQDN"
	CORTEX_API_URL_ENV_VAR                = "CORTEX_API_URL"
	CORTEX_API_KEY_ENV_VAR                = "CORTEX_API_KEY"
	CORTEX_API_KEY_ID_ENV_VAR             = "CORTEX_API_KEY_ID"
	CORTEX_API_KEY_TYPE_ENV_VAR           = "CORTEX_API_KEY_TYPE"
	CORTEX_HEADERS_ENV_VAR                = "CORTEX_HEADERS"
	CORTEX_AGENT_ENV_VAR                  = "CORTEX_AGENT"
	CORTEX_SKIP_SSL_VERIFY_ENV_VAR        = "CORTEX_SKIP_SSL_VERIFY"
	CORTEX_CONFIG_FILE_ENV_VAR            = "CORTEX_CONFIG_FILE"
	CORTEX_TIMEOUT_ENV_VAR                = "CORTEX_TIMEOUT"
	CORTEX_MAX_RETRIES_ENV_VAR            = "CORTEX_MAX_RETRIES"
	CORTEX_RETRY_MAX_DELAY_ENV_VAR        = "CORTEX_RETRY_MAX_DELAY"
	CORTEX_CRASH_STACK_DIR_ENV_VAR        = "CORTEX_CRASH_STACK_DIR"
	CORTEX_LOG_LEVEL_ENV_VAR              = "CORTEX_LOG_LEVEL"
	CORTEX_SKIP_LOGGING_TRANSPORT_ENV_VAR = "CORTEX_SKIP_LOGGING_TRANSPORT"
)

type Config struct {
	checkEnvironmentVars bool
	cortexFQDN           string            `json:"fqdn"`
	cortexAPIURL         string            `json:"api_url"`
	cortexAPIKey         string            `json:"api_key"`
	cortexAPIKeyID       int               `json:"api_key_id"`
	cortexAPIKeyType     string            `json:"api_key_type"`
	headers              map[string]string `json:"headers"`
	agent                string            `json:"agent"`
	skipSSLVerify        bool              `json:"skip_ssl_verify"`
	transport            *http.Transport   `json:"-"`
	timeout              int               `json:"timeout"`
	maxRetries           int               `json:"max_retries"`
	retryMaxDelay        int               `json:"retry_max_delay"`
	crashStackDir        string            `json:"crash_stack_dir"`
	logLevel             string            `json:"log_level"`
	logger               cortexLog.Logger  `json:"-"`
	skipLoggingTransport bool              `json:"skip_logging_transport"`
}

// CortexFQDN returns the FQDN of the Cortex tenant.
func (c *Config) CortexFQDN() string { return c.cortexFQDN }

// CortexAPIURL returns the API URL for the Cortex.
func (c *Config) CortexAPIURL() string { return c.cortexAPIURL }

// CortexAPIKeyType returns the Cortex API key type.
func (c *Config) CortexAPIKeyType() string { return c.cortexAPIKeyType }

// CortexAPIKey returns the Cortex API key.
func (c *Config) CortexAPIKey() string { return c.cortexAPIKey }

// CortexAPIKeyID returns the Cortex API key ID.
func (c *Config) CortexAPIKeyID() int { return c.cortexAPIKeyID }

// Headers returns the HTTP headers.
func (c *Config) Headers() map[string]string { return c.headers }

// Agent returns the user agent.
func (c *Config) Agent() string { return c.agent }

// SkipSSLVerify returns whether to skip TLS certificate verification.
func (c *Config) SkipSSLVerify() bool { return c.skipSSLVerify }

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
		checkEnvironmentVars: true,
		cortexFQDN:           "",
		cortexAPIURL:         "",
		cortexAPIKeyType:     "advanced",
		headers:              make(map[string]string),
		agent:                "",
		skipSSLVerify:        false,
		transport:            http.DefaultTransport.(*http.Transport),
		timeout:              30, // 30 seconds
		maxRetries:           3,
		retryMaxDelay:        60, // 60 seconds
		crashStackDir:        os.TempDir(),
		logLevel:             "info",
		logger:               nil,
		skipLoggingTransport: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.checkEnvironmentVars {
		config.overwriteFromEnvVars()
	}

	// Populate API URL using FQDN if no value is configured
	//if config.cortexAPIURL == "" || !strings.HasPrefix(config.cortexAPIURL, "https://api-") {
	if config.cortexAPIURL == "" {
		config.cortexAPIURL = fqdnToAPIURL(config.cortexFQDN)
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

	if cFile.cortexAPIURL == "" || !strings.HasPrefix(cFile.cortexAPIURL, "https://api-") {
		cFile.cortexAPIURL = fqdnToAPIURL(cFile.cortexFQDN)
	}
	return NewConfig(
		WithCortexFQDN(cFile.cortexFQDN),
		WithCortexAPIURL(cFile.cortexAPIURL),
		WithCortexAPIKey(cFile.cortexAPIKey),
		WithCortexAPIKeyID(cFile.cortexAPIKeyID),
		WithCortexAPIKeyType(cFile.cortexAPIKeyType),
		WithCheckEnvironment(checkEnvironment),
		WithHeaders(cFile.headers),
		WithAgent(cFile.agent),
		WithSkipSSLVerify(cFile.skipSSLVerify),
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
		WithCortexFQDN(c.cortexFQDN),
		WithCortexAPIURL(c.cortexAPIURL),
		WithCortexAPIKey(c.cortexAPIKey),
		WithCortexAPIKeyID(c.cortexAPIKeyID),
		WithCortexAPIKeyType(c.cortexAPIKeyType),
		WithHeaders(c.headers),
		WithAgent(c.agent),
		WithSkipSSLVerify(c.skipSSLVerify),
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
	// - Ensure `api-` prefix
	// - healthcheck
	// - API key type
	if c.cortexFQDN == "" && c.cortexAPIURL == "" {
		return fmt.Errorf("must define at least one of the FQDN and API URL fields")
	}
	if c.cortexFQDN != "" && c.cortexAPIURL == "" {
		return fmt.Errorf("API URL value not set")
	}

	return nil
}

// SetDefaults sets default values for the configuration.
func (c *Config) SetDefaults() {
	if c.logger == nil {
		c.logger = cortexLog.DefaultLogger{Logger: log.Default()}
	}
	if c.cortexAPIKeyType == "" {
		c.cortexAPIKeyType = "standard"
	}
}

func (c *Config) overwriteFromEnvVars() {
	if envFQDN, ok := os.LookupEnv(CORTEX_FQDN_ENV_VAR); ok {
		c.cortexFQDN = envFQDN
	}
	if envAPIURL, ok := os.LookupEnv(CORTEX_API_URL_ENV_VAR); ok {
		c.cortexAPIURL = envAPIURL
	}

	if envAPIKey, ok := os.LookupEnv(CORTEX_API_KEY_ENV_VAR); ok {
		c.cortexAPIKey = envAPIKey
	}

	if envAPIKeyID, ok := os.LookupEnv(CORTEX_API_KEY_ID_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envAPIKeyID); err == nil {
			c.cortexAPIKeyID = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEX_API_KEY_ID_ENV_VAR, envAPIKeyID)
		}
	}

	if envAPIKeyType, ok := os.LookupEnv(CORTEX_API_KEY_TYPE_ENV_VAR); ok {
		c.cortexAPIKeyType = envAPIKeyType
	}

	if envHeaders, ok := os.LookupEnv(CORTEX_HEADERS_ENV_VAR); ok {
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

	if envAgent, ok := os.LookupEnv(CORTEX_AGENT_ENV_VAR); ok {
		c.agent = envAgent
	}

	if envSkipSSLVerify, ok := os.LookupEnv(CORTEX_SKIP_SSL_VERIFY_ENV_VAR); ok {
		if parsedBool, err := strconv.ParseBool(envSkipSSLVerify); err == nil {
			c.skipSSLVerify = parsedBool
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected true/false.\n", CORTEX_SKIP_SSL_VERIFY_ENV_VAR, envSkipSSLVerify)
		}
	}

	if envTimeout, ok := os.LookupEnv(CORTEX_TIMEOUT_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envTimeout); err == nil {
			c.timeout = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEX_TIMEOUT_ENV_VAR, envTimeout)
		}
	}

	if envMaxRetries, ok := os.LookupEnv(CORTEX_MAX_RETRIES_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envMaxRetries); err == nil {
			c.maxRetries = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEX_MAX_RETRIES_ENV_VAR, envMaxRetries)
		}
	}

	if envRetryMaxDelay, ok := os.LookupEnv(CORTEX_RETRY_MAX_DELAY_ENV_VAR); ok {
		if parsedInt, err := strconv.Atoi(envRetryMaxDelay); err == nil {
			c.retryMaxDelay = parsedInt
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected integer.\n", CORTEX_RETRY_MAX_DELAY_ENV_VAR, envRetryMaxDelay)
		}
	}

	if envCrashStackDir, ok := os.LookupEnv(CORTEX_CRASH_STACK_DIR_ENV_VAR); ok {
		c.crashStackDir = envCrashStackDir
	}

	if envLogLevel, ok := os.LookupEnv(CORTEX_LOG_LEVEL_ENV_VAR); ok {
		c.logLevel = envLogLevel
	}

	if envSkipLoggingTransport, ok := os.LookupEnv(CORTEX_SKIP_LOGGING_TRANSPORT_ENV_VAR); ok {
		if parsedBool, err := strconv.ParseBool(envSkipLoggingTransport); err == nil {
			c.skipLoggingTransport = parsedBool
		} else {
			fmt.Printf("Warning: Invalid value for %s environment variable: %s. Expected true/false.\n", CORTEX_SKIP_LOGGING_TRANSPORT_ENV_VAR, envSkipLoggingTransport)
		}
	}
}

func fqdnToAPIURL(fqdn string) string {
	fqdnLowercase := strings.ToLower(fqdn)
	if strings.HasPrefix(fqdnLowercase, "https://api-") {
		return fqdn
	}

	ensureAPIPrefix := func(hostname string) string {
		if !strings.HasPrefix(hostname, "api-") {
			return "api-" + hostname
		}
		return hostname
	}

	if !strings.HasPrefix(strings.ToLower(fqdn), "https://") {
		fqdnParts := strings.SplitN(fqdn, "://", 2)
		if len(fqdnParts) != 2 {
			return "https://" + ensureAPIPrefix(fqdn)
		}
		return "https://" + ensureAPIPrefix(fqdnParts[1])
	}

	return ensureAPIPrefix(fqdn)
}
