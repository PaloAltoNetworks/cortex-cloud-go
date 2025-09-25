// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
)

// API endpoint path specification.
const (
	ListUsersEndpoint      = "public_api/v1/rbac/get_users/"
	ListRolesEndpoint      = "public_api/v1/rbac/get_roles/"
	SetUserRoleEndpoint    = "public_api/v1/rbac/set_user_role/"
	GetRiskScoreEndpoint   = "public_api/v1/get_risk_score/"
	ListRiskyUsersEndpoint = "public_api/v1/risk/get_risky_users/"
	ListRiskyHostsEndpoint = "public_api/v1/risky_hosts/"

	// Asset Group Endpoints
	CreateAssetGroupEndpoint = "public_api/v1/asset-groups/create/"
	UpdateAssetGroupEndpoint = "public_api/v1/asset-groups/update/"
	DeleteAssetGroupEndpoint = "public_api/v1/asset-groups/delete/"
	ListAssetGroupsEndpoint  = "public_api/v1/asset-groups"

	// Auth Settings Endpoints
	ListIDPMetadataEndpoint  = "public_api/v1/sso/get_idp_metadata/"
	ListAuthSettingsEndpoint = "public_api/v1/sso/get_sso_config/"
	CreateAuthSettingsEndpoint = "public_api/v1/sso/set_config/"
	UpdateAuthSettingsEndpoint = "public_api/v1/sso/set_config/"
	DeleteAuthSettingsEndpoint = "public_api/v1/sso/delete_config/"
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
