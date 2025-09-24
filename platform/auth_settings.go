// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform/types"
)

// ListIDPMetadata returns the metadata for all IDPs.
//
// This endpoint requires Instance Administrator permissions.
func (c *Client) ListIDPMetadata(ctx context.Context) (types.ListIDPMetadataResponse, error) {
	var ans types.ListIDPMetadataResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListIDPMetadataEndpoint, nil, nil, types.ListIDPMetadataRequestData{}, &ans, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})
	return ans, err
}

// ListAuthSettings returns the authentication settings for all configured
// domains in the tenant.
//
// This endpoint requires Instance Administrator permissions.
func (c *Client) ListAuthSettings(ctx context.Context) ([]types.AuthSettings, error) {
	var ans []types.AuthSettings
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListAuthSettingsEndpoint, nil, nil, types.ListAuthSettingsRequestData{}, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

// CreateAuthSettings creates authentication settings for the specified domain
// using the provided IDP SSO or metadata URL.
//
// To configure IDP SSO, the `idp_sso_url`, `idp_issuer` and `idp_certificate`
// fields are required. To configure via metadata URL, the `metadata_url` is
// the only required field.
//
// This endpoint requires Instance Administrator permissions.
func (c *Client) CreateAuthSettings(ctx context.Context, req types.CreateAuthSettingsRequestData) (bool, error) {
	var resp bool
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateAuthSettingsEndpoint, nil, nil, req, &resp, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return resp, err
}

// UpdateAuthSettings updates the existing authentication settings for the
// specified domain.
//
// To update the default domain, provide empty strings for the
// `current_domain_value` and `new_domain_value` fields.
//
// This endpoint requires Instance Administrator permissions.
func (c *Client) UpdateAuthSettings(ctx context.Context, req types.UpdateAuthSettingsRequestData) (bool, error) {
	var resp bool
	_, err := c.internalClient.Do(ctx, http.MethodPost, UpdateAuthSettingsEndpoint, nil, nil, req, &resp, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return resp, err
}

// DeleteAuthSettings deletes all authentication settings for the specified
// domain.
//
// This endpoint requires Instance Administrator permissions.
func (c *Client) DeleteAuthSettings(ctx context.Context, domain string) (bool, error) {
	var resp bool
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteAuthSettingsEndpoint, nil, nil, types.DeleteAuthSettingsRequestData{ Domain: domain, }, &resp, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return resp, err
}
