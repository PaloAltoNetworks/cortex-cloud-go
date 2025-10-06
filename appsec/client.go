// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package appsec

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
)

// API endpoint path specification.
const (
	RulesEndpoint           = "public_api/appsec/v1/rules"
	RulesValidationEndpoint = "public_api/appsec/v1/rules/validate"
	RulesActionsEndpoint    = "public_api/appsec/v1/rules/rule_actions"
)

// Client is the client for the namespace.
type Client struct {
	internalClient *client.Client
}

// NewClient returns a new client for this namespace.
func NewClient(config *client.Config) (*Client, error) {
	internalClient, err := client.NewClientFromConfig(config)
	return &Client{internalClient: internalClient}, err
}

func (c *Client) BuildInfo() map[string]string {
	return map[string]string{
		"gitCommit": GitCommit,
		"goVersion": GoVersion,
		"buildDate": BuildDate,
	}
}
