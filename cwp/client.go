// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
)

const (
	ListCloudWorkloadPoliciesEndpoint     = "public_api/v1/cwp/get_policies"
	CreateCloudWorkloadPolicyEndpoint     = "public_api/v1/cwp/create_policy"
	GetCloudWorkloadPolicyDetailsEndpoint = "public_api/v1/cwp/get_policy_details"
	EditCloudWorkloadPolicyEndpoint       = "public_api/v1/cwp/edit_policy"
	DeleteCloudWorkloadPolicyEndpoint     = "public_api/v1/cwp/delete_policy"
)

// Client is the client for the namespace.
type Client struct {
	internalClient *app.Client
}

// NewClient returns a new client for this namespace.
func NewClient(config *api.Config) (*Client, error) {
	internalClient, err := app.NewClient(config)
	return &Client{internalClient: internalClient}, err
}
