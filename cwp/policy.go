// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"dario.cat/mergo"
	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types/cwp"
	convert "github.com/PaloAltoNetworks/cortex-cloud-go/types/util"
)

// Add Policy
func (c *Client) CreatePolicy(ctx context.Context, input types.CreatePolicyRequest) (types.CreatePolicyResponse, error) {
	var res types.CreatePolicyResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreatePolicyEndpoint, nil, nil, input, &res, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return res, err
}

// Get Policy By ID
func (c *Client) GetPolicyByID(ctx context.Context, policyId string) (types.Policy, error) {
	var res types.Policy
	_, err := c.internalClient.Do(ctx, http.MethodGet, GetPolicyByIDEndpoint, &[]string{policyId}, nil, nil, &res, nil)
	return res, err
}

// List Policies
func (c *Client) ListPolicies(ctx context.Context, policyTypes []string) ([]types.Policy, error) {
	queryParams := convert.StringSliceToQuery("policy_types", policyTypes)
	var res []types.Policy
	_, err := c.internalClient.Do(ctx, http.MethodGet, ListPoliciesEndpoint, nil, &queryParams, nil, &res, nil)
	return res, err
}

// Delete Policy
func (c *Client) DeletePolicy(ctx context.Context, policyID int, closeIssues bool) error {
	queryVals := convert.StringToQuery("close_issues", strconv.FormatBool(closeIssues))
	_, err := c.internalClient.Do(ctx, http.MethodDelete, DeletePolicyEndpoint, &[]string{strconv.Itoa(policyID)}, &queryVals, nil, nil, nil)
	return err
}

// Update Policy
func (c *Client) UpdatePolicy(ctx context.Context, input types.UpdatePolicyRequest) error {
	policy, err := c.GetPolicyByID(ctx, input.Id)
	if err != nil {
		return err
	}
	if err = mergo.Merge(&input, policy); err != nil {
		return fmt.Errorf("failed to merge policies during update (we expected this might happen): %s", err.Error())
	}

	_, err = c.internalClient.Do(ctx, http.MethodPut, EditPolicyEndpoint, &[]string{input.Id}, nil, policy, nil, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})

	return err
}
