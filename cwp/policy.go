// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package cwp provides cwp
package cwp

import (
	"context"
	"net/http"
	"strconv"

	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/cwp"
	convert "github.com/PaloAltoNetworks/cortex-cloud-go/types/util"
)

// CreatePolicy method
func (c *Client) CreatePolicy(ctx context.Context, input types.CreatePolicyRequest) (types.CreatePolicyResponse, error) {
	var res types.CreatePolicyResponse

	_, err := c.internalClient.Do(ctx, http.MethodPost, CreatePolicyEndpoint, nil, nil, input, &res, nil)
	if err != nil {
		c.Logger().Error(ctx, "Failed to create policy", map[string]interface{}{
			"error": err.Error(),
		})
	}

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
	queryParams := convert.StringSliceToQuery("types", policyTypes)
	var res []types.Policy
	_, err := c.internalClient.Do(ctx, http.MethodGet, ListPoliciesEndpoint, nil, &queryParams, nil, &res, nil)
	return res, err
}

// Delete Policy - original function kept for backward compatibility
func (c *Client) DeletePolicy(ctx context.Context, policyID int, closeIssues bool) error {
	queryVals := convert.StringToQuery("close_issues", strconv.FormatBool(closeIssues))
	_, err := c.internalClient.Do(ctx, http.MethodDelete, DeletePolicyEndpoint, &[]string{strconv.Itoa(policyID)}, &queryVals, nil, nil, nil)
	return err
}

// DeletePolicyByString deletes a policy using a string ID (for UUIDs)
func (c *Client) DeletePolicyByString(ctx context.Context, policyID string, closeIssues bool) error {
	queryVals := convert.StringToQuery("close_issues", strconv.FormatBool(closeIssues))
	_, err := c.internalClient.Do(ctx, http.MethodDelete, DeletePolicyEndpoint, &[]string{policyID}, &queryVals, nil, nil, nil)
	return err
}

// UpdatePolicy method updates an existing policy
func (c *Client) UpdatePolicy(ctx context.Context, input types.UpdatePolicyRequest) error {
	// Get the existing policy to merge with input
	existingPolicy, err := c.GetPolicyByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// Create a new UpdatePolicyRequest based on the existing policy
	// This ensures all fields from the existing policy are present
	mergedInput := types.UpdatePolicyRequest{
		Policy: existingPolicy,
	}

	// Now, selectively override fields in mergedInput with values from the 'input'
	// only if the 'input' provides a non-zero/non-empty value for that field.
	copyNonZeroFields(&mergedInput, input)

	// Log the merged request for debugging
	c.Logger().Debug(ctx, "Updating policy with merged input", map[string]interface{}{
		"input": mergedInput,
	})

	// Make the actual update API call with merged data
	_, err = c.internalClient.Do(ctx, http.MethodPut, CreatePolicyEndpoint, &[]string{mergedInput.ID}, nil, mergedInput, nil, nil)
	if err != nil {
		c.Logger().Error(ctx, "Failed to update policy", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return err
}

// copyNonZeroFields copies non-zero/non-empty fields from src to dest.
// This is a generic helper that can be extended or made more specific if needed.
func copyNonZeroFields(dest *types.UpdatePolicyRequest, src types.UpdatePolicyRequest) {
	if src.Name != "" {
		dest.Name = src.Name
	}
	if src.Type != "" {
		dest.Type = src.Type
	}
	if src.Description != "" {
		dest.Description = src.Description
	}
	if len(src.EvaluationModes) > 0 {
		dest.EvaluationModes = src.EvaluationModes
	}
	if src.EvaluationStage != "" {
		dest.EvaluationStage = src.EvaluationStage
	}
	if len(src.RulesIDs) > 0 {
		dest.RulesIDs = src.RulesIDs
	}
	if src.Condition != "" {
		dest.Condition = src.Condition
	}
	if src.Exception != "" {
		dest.Exception = src.Exception
	}
	if src.AssetScope != "" {
		dest.AssetScope = src.AssetScope
	}
	if len(src.AssetGroupIDs) > 0 {
		dest.AssetGroupIDs = src.AssetGroupIDs
	}
	if len(src.AssetGroups) > 0 {
		dest.AssetGroups = src.AssetGroups
	}
	if src.PolicyAction != "" {
		dest.PolicyAction = src.PolicyAction
	}
	if src.PolicySeverity != "" {
		dest.PolicySeverity = src.PolicySeverity
	}
	if src.RemediationGuidance != "" {
		dest.RemediationGuidance = src.RemediationGuidance
	}
	// Boolean field needs special handling - we always want to take the input value
	// as it represents the desired state, even if it's 'false'.
	// The original code already handled this correctly.
	dest.Disabled = src.Disabled
}
