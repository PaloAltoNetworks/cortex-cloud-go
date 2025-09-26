// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
)

const (
	ListPoliciesEndpoint     = "public_api/v1/cwp/get_policies"
	CreatePolicyEndpoint     = "public_api/v1/cwp/policies"
	GetPolicyByIDEndpoint = "public_api/v1/cwp/get_policy_details"
	EditPolicyEndpoint       = "public_api/v1/cwp/edit_policy"
	DeletePolicyEndpoint     = "public_api/v1/cwp/delete_policy"
)

// Client is the client for the namespace.
type Client struct {
	internalClient *client.Client
}

// NewClient returns a new client for this namespace.
func NewClient(config *api.Config) (*Client, error) {
	internalClient, err := client.NewClient(config)
	return &Client{internalClient: internalClient}, err
}
