// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
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

// Client is the client for the namespace.
type Client struct {
	internalClient *client.Client
}

// NewClient returns a new client for this namespace.
func NewClient(config *client.Config) (*Client, error) {
	internalClient, err := client.NewClientFromConfig(config)
	return &Client{internalClient: internalClient}, err
}
