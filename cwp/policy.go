// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"fmt"
	"context"
	"net/http"
	"strconv"

	"dario.cat/mergo"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/util"
	"github.com/PaloAltoNetworks/cortex-cloud-go/cwp/types"
)

// Add Policy
func (c *Client) CreatePolicy(ctx context.Context, input types.CreatePolicyRequest) (types.CreatePolicyResponse, error) {
	var res types.CreatePolicyResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreatePolicyEndpoint, nil, nil, input, &res, nil)
	return res, err
}

// Get Policy By ID
func (c *Client) GetPolicyByID(ctx context.Context, policyId string) (types.Policy, error) {
	//var res types.GetPolicyByIDResponse
	var res types.Policy
	_, err := c.internalClient.Do(ctx, http.MethodGet, GetPolicyByIDEndpoint, &[]string{policyId}, nil, nil, &res, nil)
	return res, err
}

// List Policies
//func (c *Client) ListPolicies(ctx context.Context, input types.ListPoliciesRequest) (types.ListPoliciesResponse, error) {
func (c *Client) ListPolicies(ctx context.Context, policyTypes []string) ([]types.Policy, error) {
	//var res types.ListPoliciesResponse
	queryParams := util.StringSliceToQuery("policy_types", policyTypes)
	var res []types.Policy
	_, err := c.internalClient.Do(ctx, http.MethodGet, ListPoliciesEndpoint, nil, &queryParams, nil, &res, nil)
	return res, err
}

// Delete Policy
//func (c *Client) DeletePolicy(ctx context.Context, input types.DeletePolicyRequest) error {
func (c *Client) DeletePolicy(ctx context.Context, policyID int, closeIssues bool) error {
	//id := strconv.Itoa(policyID)
	//req := types.DeletePolicyRequest{
	//	//Id: id,
	//	CloseIssues: closeIssues,
	//}
	//deleteQueryValues := req.ToQueryValues()
	deleteQueryValues := util.StringToQuery("close_issues", strconv.FormatBool(closeIssues))
	_, err := c.internalClient.Do(ctx, http.MethodDelete, DeletePolicyEndpoint, &[]string{ strconv.Itoa(policyID) }, &deleteQueryValues, nil, nil, nil)
	return err
}

// Update Policy
func (c *Client) UpdatePolicy(ctx context.Context, input types.UpdatePolicyRequest) error {
	policy, err := c.GetPolicyByID(ctx, input.Id)
	if err != nil {
		return err
	}

	//updatedPolicy := policy.Policy.ToUpdateRequest()
	//updatedPolicy := policy.Policy.ToUpdateRequest()

	//if err = mergo.Merge(&updatedPolicy, input); err != nil {
	if err = mergo.Merge(&input, policy); err != nil {
		//return err
		return fmt.Errorf("failed to merge policies during update (we expected this might happen): %s", err.Error())
	}

	//_, err = c.internalClient.Do(ctx, http.MethodPut, EditPolicyEndpoint, &[]string{input.Id}, nil, updatedPolicy, nil)
	_, err = c.internalClient.Do(ctx, http.MethodPut, EditPolicyEndpoint, &[]string{input.Id}, nil, policy, nil, nil)

	return err
}
