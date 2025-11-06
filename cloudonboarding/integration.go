// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types/cloudonboarding"
)

// CreateTemplate creates a new Cloud Onboarding Integration Template.
//
// TODO: details
func (c *Client) CreateIntegrationTemplate(ctx context.Context, input *types.CreateIntegrationTemplateRequest) (types.CreateTemplateOrEditIntegrationInstanceResponse, error) {
	var ans types.CreateTemplateOrEditIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateIntegrationTemplateEndpoint, nil, nil, input, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// TODO: Create alternate version that doesn't marshal the response into IntegrationInstance
// GetDetails returns the configuration details of the specified integration instance.
func (c *Client) GetIntegrationInstanceDetails(ctx context.Context, instanceID string) (types.IntegrationInstance, error) {
	var ans types.GetIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetIntegrationInstanceDetailsEndpoint, nil, nil, types.NewGetIntegrationInstanceRequest(instanceID), &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})
	if err != nil {
		return types.IntegrationInstance{}, err
	}
	return ans.Marshal()
}

// TODO: Create alternate version that doesn't marshal the response into IntegrationInstance
func (c *Client) ListIntegrationInstances(ctx context.Context, input *types.ListIntegrationInstancesRequest) ([]types.IntegrationInstance, error) {
	var ans types.ListIntegrationInstancesResponseWrapper
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListIntegrationInstancesEndpoint, nil, nil, input, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})
	if err != nil {
		return []types.IntegrationInstance{}, err
	}
	return ans.Marshal()
}
func (c *Client) EditIntegrationInstance(ctx context.Context, input *types.EditIntegrationInstanceRequest) (types.CreateTemplateOrEditIntegrationInstanceResponse, error) {
	var ans types.CreateTemplateOrEditIntegrationInstanceResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, EditIntegrationInstanceEndpoint, nil, nil, input, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})

	return ans, err
}

func (c *Client) EnableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, types.NewEnableOrDisableIntegrationInstancesRequest(instanceIDs, true), nil, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return err
}

func (c *Client) DisableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, types.NewEnableOrDisableIntegrationInstancesRequest(instanceIDs, false), nil, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return err
}

func (c *Client) DeleteIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteIntegrationInstancesEndpoint, nil, nil, types.NewDeleteIntegrationInstanceRequest(instanceIDs), nil, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return err
}
