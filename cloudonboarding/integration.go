// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
	"github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding/types"
)

// CreateTemplate creates a new Cloud Onboarding Integration Template.
//
// TODO: details
func (c *Client) CreateIntegrationTemplate(ctx context.Context, input types.CreateIntegrationTemplateRequest) (types.CreateTemplateOrEditIntegrationInstanceResponse, error) {
	var ans types.CreateTemplateOrEditIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateIntegrationTemplateEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

// GetDetails returns the configuration details of the specified integration instance.
func (c *Client) GetIntegrationInstanceDetails(ctx context.Context, input types.GetIntegrationInstanceRequest) (types.GetIntegrationInstanceResponse, error) {
	var ans types.GetIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetIntegrationInstanceDetailsEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

func (c *Client) ListIntegrationInstances(ctx context.Context, input types.ListIntegrationInstancesRequest) (types.ListIntegrationInstancesResponseWrapper, error) {
	var ans types.ListIntegrationInstancesResponseWrapper
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListIntegrationInstancesEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}
func (c *Client) EditIntegrationInstance(ctx context.Context, input types.EditIntegrationInstanceRequest) (types.CreateTemplateOrEditIntegrationInstanceResponse, error) {
	var ans types.CreateTemplateOrEditIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, EditIntegrationInstanceEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

func (c *Client) EnableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, types.EnableOrDisableIntegrationInstancesRequest{ IDs: instanceIDs, Enable: true }, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})
	return err
}

func (c *Client) DisableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, types.EnableOrDisableIntegrationInstancesRequest{ IDs: instanceIDs, Enable: false }, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})
	return err
}

func (c *Client) DeleteIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteIntegrationInstancesEndpoint, nil, nil, types.DeleteIntegrationInstanceRequest{ IDs: instanceIDs }, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})
	return err
}
