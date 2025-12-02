// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package cwp provides cwp
package cwp

import (
	"context"
	"net/http"
	"strconv"

	"dario.cat/mergo"

	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/cwp"
	convert "github.com/PaloAltoNetworks/cortex-cloud-go/types/util"
)

// CreateCloudWorkloadPolicy creates a new Cloud Workload Policy.
func (c *Client) CreateCloudWorkloadPolicy(ctx context.Context, input types.CreateCloudWorkloadPolicyRequest) (policyID string, err error) {
	var res types.CreateCloudWorkloadPolicyResponse
	_, err = c.internalClient.Do(ctx, http.MethodPost, CreatePolicyEndpoint, nil, nil, input, &res, nil)
	if err != nil {
		return "", err
	}
	return res.ID, err
}

// GetCloudWorkloadPolicyByID retrieves the Cloud Workload Policy for the given
// ID.
func (c *Client) GetCloudWorkloadPolicyByID(ctx context.Context, policyID string) (policy *types.CloudWorkloadPolicy, err error) {
	_, err = c.internalClient.Do(ctx, http.MethodGet, GetPolicyByIDEndpoint, &[]string{policyID}, nil, nil, policy, nil)
	if err != nil {
		return nil, err
	}
	return policy, err
}

// ListCloudWorkloadPolicies retrieves all Cloud Workload Policies with the
// specified policy types.
//
// If no policy types are specified, all policies will be returned.
func (c *Client) ListCloudWorkloadPolicies(ctx context.Context, policyTypes []string) (policies *[]types.CloudWorkloadPolicy, err error) {
	queryValues := convert.StringSliceToQuery("types", policyTypes)
	_, err = c.internalClient.Do(ctx, http.MethodGet, ListPoliciesEndpoint, nil, &queryValues, nil, policies, nil)
	if err != nil {
		return nil, err
	}
	return policies, err
}

// DeletePolicyByString deletes the Cloud Workload Policy with the given ID.
func (c *Client) DeleteCloudWorkloadPolicy(ctx context.Context, policyID string, closeIssues bool) (success bool, err error) {
	queryValues := convert.StringToQuery("close_issues", strconv.FormatBool(closeIssues))
	_, err = c.internalClient.Do(ctx, http.MethodDelete, DeletePolicyEndpoint, &[]string{policyID}, &queryValues, nil, nil, nil)
	return err != nil, err
}

// UpdatePolicy method updates an existing policy
func (c *Client) UpdatePolicy(ctx context.Context, input types.UpdateCloudWorkloadPolicyRequest) (success bool, err error) {
	// Fetch policy from API
	existingPolicy, err := c.GetCloudWorkloadPolicyByID(ctx, input.ID)
	if err != nil {
		return false, err
	}

	// Merge with input, giving preference to input values
	var merged types.UpdateCloudWorkloadPolicyRequest
	if err = mergo.Merge(&merged, existingPolicy, mergo.WithOverride); err != nil {
		return false, err
	}

	_, doErr := c.internalClient.Do(ctx, http.MethodPut, CreatePolicyEndpoint, &[]string{merged.ID}, nil, merged, nil, nil)
	return doErr != nil, doErr
}
