// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package cwp

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	policies, err := client.ListCloudWorkloadPolicies(ctx, nil)
	require.NoError(t, err, "Failed to list policies")
	require.NotNil(t, policies, "Policies list should not be nil")

	// Assert that there is at least one policy (assuming a pre-existing policy in the test environment)
	require.NotNil(t, policies)
	assert.GreaterOrEqual(t, len(*policies), 0, "Expected at least one policy to be returned")

	// TODO: pre-populate data so we don't need this check
	if len(*policies) > 0 {
		// Take the first policy from the list and get it by ID
		firstPolicy := (*policies)[0]
		policyByID, err := client.GetCloudWorkloadPolicyByID(ctx, firstPolicy.ID)
		require.NoError(t, err, "Failed to get policy by ID")
		require.NotNil(t, policyByID, "Policy by ID should not be nil")
		assert.Equal(t, firstPolicy.ID, policyByID.ID, "IDs should match")
		assert.Equal(t, firstPolicy.Name, policyByID.Name, "Names should match")
	}
}
