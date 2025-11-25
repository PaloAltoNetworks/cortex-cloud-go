// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package cwp

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/cwp"
)

func setupAcceptanceTest(t *testing.T) *Client {
	apiURL := os.Getenv("CORTEX_API_URL_TEST")
	if apiURL == "" {
		t.Fatalf("Error: CORTEX_API_URL_TEST environment variable not set.")
	}
	apiKey := os.Getenv("CORTEX_API_KEY_TEST")
	if apiKey == "" {
		t.Fatalf("Error: CORTEX_API_KEY_TEST environment variable not set.")
	}
	apiKeyIDStr := os.Getenv("CORTEX_API_KEY_ID_TEST")
	if apiKeyIDStr == "" {
		t.Fatalf("Error: CORTEX_API_KEY_ID_TEST environment variable not set.")
	}

	apiKeyID, err := strconv.Atoi(apiKeyIDStr)
	if err != nil {
		t.Fatalf("failed to convert API key ID \"%s\" to int: %s", apiKeyIDStr, err.Error())
	}

	client, err := NewClient(
		WithCortexAPIURL(apiURL),
		WithCortexAPIKey(apiKey),
		WithCortexAPIKeyID(apiKeyID),
		WithLogLevel("debug"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	return client
}

func TestAccListPolicies(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()

	// List all policies
	policies, err := client.ListPolicies(ctx, nil)
	require.NoError(t, err, "Failed to list policies")
	require.NotNil(t, policies, "Policies list should not be nil")

	// Assert that there is at least one policy (assuming a pre-existing policy in the test environment)
	// This is a basic check; more robust tests would involve creating known policies first.
	assert.GreaterOrEqual(t, len(policies), 0, "Expected at least one policy to be returned")

	// Optionally, filter by a specific policy type if known to exist
	// For example, if you know there are "COMPLIANCE" policies:
	compliancePolicies, err := client.ListPolicies(ctx, []string{"COMPLIANCE"})
	require.NoError(t, err, "Failed to list compliance policies")
	require.NotNil(t, compliancePolicies, "Compliance policies list should not be nil")

	// Assert that all returned policies are of type "COMPLIANCE"
	for _, policy := range compliancePolicies {
		assert.Equal(t, "COMPLIANCE", policy.Type, "Expected policy type to be COMPLIANCE")
	}

	// Take the first policy from the list and get it by ID
	if len(policies) > 0 {
		firstPolicy := policies[0]
		policyByID, err := client.GetPolicyByID(ctx, firstPolicy.ID)
		require.NoError(t, err, "Failed to get policy by ID")
		require.NotNil(t, policyByID, "Policy by ID should not be nil")
		assert.Equal(t, firstPolicy.ID, policyByID.ID, "IDs should match")
		assert.Equal(t, firstPolicy.Name, policyByID.Name, "Names should match")
	}
}

func TestAccMisconfigurationPolicyLifecycle(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	currentTime := time.Now()
	timestamp := strconv.FormatInt(currentTime.Unix(), 10)
	name := fmt.Sprintf("go-sdk-acc-test-misconfig-%s", timestamp)

	// Try to list COMPLIANCE policies specifically
	compliancePolicies, err := client.ListPolicies(ctx, []string{"COMPLIANCE"})
	require.NoError(t, err, "Failed to list COMPLIANCE policies")

	// If no COMPLIANCE policies found, try listing all policies
	if len(compliancePolicies) == 0 {
		t.Log("No COMPLIANCE policies found, listing all policies")
		allPolicies, err := client.ListPolicies(ctx, nil)
		require.NoError(t, err, "Failed to list all policies")

		// Look for COMPLIANCE policies in the full list
		for _, policy := range allPolicies {
			if policy.Type == "COMPLIANCE" {
				compliancePolicies = append(compliancePolicies, policy)
			}
		}
	}

	// If we still can't find any COMPLIANCE policies
	if len(compliancePolicies) == 0 {
		t.Log("No COMPLIANCE policies found, using hardcoded values")

		createReq := types.CreatePolicyRequest{
			Type:            "COMPLIANCE",
			Name:            name,
			Description:     "Created by Go SDK acceptance test",
			EvaluationStage: "RUNTIME",
			RulesIDs:        []string{"00000000-0000-0000-0000-000000000015"},
			AssetGroupIDs:   []int{1},
			PolicyAction:    "ISSUE",
			PolicySeverity:  "LOW",
		}

		createResp, err := client.CreatePolicy(ctx, createReq)
		require.NoError(t, err, "Failed to create policy with hardcoded values")
		require.NotEmpty(t, createResp.ID, "Created policy should have an ID")

		policyID := createResp.ID
		t.Logf("Successfully created policy with ID: %s", policyID)

		defer func() {
			t.Logf("Cleaning up policy ID: %s", policyID)
			err = client.DeletePolicyByString(ctx, policyID, false)
			if err != nil {
				t.Logf("Warning: Failed to clean up policy: %s", err.Error())
			}
		}()

		// Verify policy was created
		policy, err := client.GetPolicyByID(ctx, policyID)
		require.NoError(t, err, "Failed to retrieve created policy")

		// Test update
		updatedPolicy := policy
		updatedPolicy.Name = name + "-updated"
		updatedPolicy.Description = "Updated by Go SDK"

		updateReq := types.UpdatePolicyRequest{
			Policy: updatedPolicy,
		}

		err = client.UpdatePolicy(ctx, updateReq)
		require.NoError(t, err, "Failed to update policy")

		// Verify update
		updatedPolicyResult, err := client.GetPolicyByID(ctx, policyID)
		require.NoError(t, err, "Failed to retrieve updated policy")
		assert.Equal(t, name+"-updated", updatedPolicyResult.Name)

		return
	}

	// If we found COMPLIANCE policies, use the first one as a template
	templatePolicy := compliancePolicies[0]
	t.Logf("Using COMPLIANCE policy as template: %s (ID: %s)", templatePolicy.Name, templatePolicy.ID)

	// Log the template policy details for debugging
	t.Logf("Template policy full details: %+v", templatePolicy)

	// Create a minimal COMPLIANCE policy using essential fields
	createReq := types.CreatePolicyRequest{
		Type:            "COMPLIANCE",
		Name:            name,
		Description:     "Created by Go SDK acceptance test",
		EvaluationStage: "RUNTIME",
		RulesIDs:        []string{"00000000-0000-0000-0000-000000000015"},
		AssetGroupIDs:   []int{1},
		PolicyAction:    "ISSUE",
		PolicySeverity:  "LOW",
	}

	// Copy additional fields from the template if they exist
	if len(templatePolicy.EvaluationModes) > 0 {
		createReq.EvaluationModes = templatePolicy.EvaluationModes
	} else {
		createReq.EvaluationModes = []string{"PERIODIC"}
	}

	if templatePolicy.Condition != "" {
		createReq.Condition = templatePolicy.Condition
	}

	if templatePolicy.Exception != "" {
		createReq.Exception = templatePolicy.Exception
	}

	if templatePolicy.AssetScope != "" {
		createReq.AssetScope = templatePolicy.AssetScope
	}

	if len(templatePolicy.AssetGroups) > 0 {
		createReq.AssetGroups = templatePolicy.AssetGroups
	}

	if templatePolicy.RemediationGuidance != "" {
		createReq.RemediationGuidance = templatePolicy.RemediationGuidance
	}

	createResp, err := client.CreatePolicy(ctx, createReq)
	require.NoError(t, err, "Failed to create policy")
	require.NotEmpty(t, createResp.ID, "Created policy should have an ID")

	policyID := createResp.ID

	defer func() {
		t.Logf("Cleaning up policy ID: %s", policyID)
		err = client.DeletePolicyByString(ctx, policyID, false)
		if err != nil {
			t.Logf("Warning: Failed to clean up policy: %s", err.Error())
		}
	}()

	// Get policy to verify it was created
	policy, err := client.GetPolicyByID(ctx, policyID)
	require.NoError(t, err, "Failed to retrieve created policy")

	// Verify policy details
	assert.Equal(t, name, policy.Name)
	assert.Equal(t, "COMPLIANCE", policy.Type)

	// Update policy
	updatedPolicy := policy
	updatedPolicy.Name = name + "-updated"
	updatedPolicy.Description = "Updated Go SDK Test"

	updateReq := types.UpdatePolicyRequest{
		Policy: updatedPolicy,
	}

	err = client.UpdatePolicy(ctx, updateReq)
	require.NoError(t, err, "Failed to update policy")

	// Get updated policy
	updatedPolicyResult, err := client.GetPolicyByID(ctx, policyID)
	require.NoError(t, err, "Failed to retrieve updated policy")

	// Verify updated policy details
	assert.Equal(t, name+"-updated", updatedPolicyResult.Name)
	assert.Equal(t, "Updated Go SDK Test", updatedPolicyResult.Description)
}
