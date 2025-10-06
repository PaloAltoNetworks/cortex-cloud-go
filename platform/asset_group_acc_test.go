// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package platform

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
	"github.com/stretchr/testify/assert"
)

func setupAcceptanceTest(t *testing.T) *Client {
	apiUrl := os.Getenv("CORTEX_API_URL_TEST")
	apiKey := os.Getenv("CORTEX_API_KEY_TEST")
	apiKeyIDStr := os.Getenv("CORTEX_API_KEY_ID_TEST")

	apiKeyID, err := strconv.Atoi(apiKeyIDStr)
	if err != nil {
		t.Fatalf("failed to convert API key ID \"%s\" to int: %s", apiKeyIDStr, err.Error())
	}

	config := &client.Config{
		ApiUrl:   apiUrl,
		ApiKey:   apiKey,
		ApiKeyId: apiKeyID,
		LogLevel: "debug",
	}

	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	return client
}

func TestAccDynamicAssetGroupLifecycle(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Create
	groupName := fmt.Sprintf("go-sdk-acctest-account-group-%s", timestamp)
	groupType := "Dynamic"
	groupDescription := fmt.Sprintf("Go SDK Acceptance Test %s", timestamp)
	andCriteria1Field := "xdm.asset.name"
	andCriteria1Type := "NCONTAINS"
	andCriteria1Value := "test"
	membershipPredicate := types.Filter{
		And: []*types.Filter{
			{
				SearchField: andCriteria1Field,
				SearchType:  andCriteria1Type,
				SearchValue: andCriteria1Value,
			},
		},
	}

	createReq := types.CreateOrUpdateAssetGroupRequest{
		GroupName:           groupName,
		GroupType:           groupType,
		GroupDescription:    groupDescription,
		MembershipPredicate: membershipPredicate,
	}

	createSuccess, groupID, err := client.CreateAssetGroup(ctx, createReq)
	if err != nil {
		t.Fatalf("error creating asset group: %s", err.Error())
	}
	if !createSuccess {
		t.Fatalf("asset group creation unsuccessful: %s", err.Error())
	}

	// Check
	assert.Equal(t, true, createSuccess)
	assert.Positive(t, groupID)

	// Defer Delete check
	testDeleteAssetGroup := func(t *testing.T, ctx context.Context, groupID int) {
		// Delete
		success, err := client.DeleteAssetGroup(ctx, groupID)

		// Check
		if err != nil {
			t.Fatalf("error deleting asset group: %s", err.Error())
		}
		if !success {
			t.Fatalf("asset group deletion unsuccessful: %s", err.Error())
		}
		assert.Equal(t, true, success)
	}
	defer testDeleteAssetGroup(t, ctx, groupID)

	// Read
	listReq := types.ListAssetGroupsRequest{
		Filters: types.Filter{
			And: []*types.Filter{
				{
					SearchField: "XDM.ASSET_GROUP.NAME",
					SearchType:  "CONTAINS",
					SearchValue: groupName,
				},
			},
		},
	}
	assetGroups, err := client.ListAssetGroups(ctx, listReq)
	if err != nil {
		t.Fatalf("error fetching asset group: %s", err.Error())
	}

	// Check
	expectedFilter := []types.AssetGroupFilter{
		{
			PrettyName: "name",
			DataType:   "TEXT",
			RenderType: "attribute",
		},
		{
			PrettyName: "doesn't contain",
			RenderType: "operator",
		},
		{
			PrettyName: andCriteria1Value,
			RenderType: "value",
		},
	}

	assert.NotNil(t, assetGroups)
	assert.Len(t, assetGroups, 1)
	assetGroup := assetGroups[0]
	assert.Equal(t, groupID, assetGroup.ID)
	assert.Equal(t, groupName, assetGroup.Name)
	assert.Equal(t, groupType, assetGroup.Type)
	assert.Equal(t, groupDescription, assetGroup.Description)
	assert.Equal(t, membershipPredicate, assetGroup.MembershipPredicate)
	assert.Equal(t, expectedFilter, assetGroup.Filter)

	// Update
	updatedGroupName := fmt.Sprintf("%s-updated", groupName)
	updatedGroupDescription := fmt.Sprintf("%s (Updated)", groupDescription)
	andCriteria2Field := "xdm.asset.name"
	andCriteria2Type := "NCONTAINS"
	andCriteria2Value := "test"
	membershipPredicate.And = append(membershipPredicate.And, &types.Filter{
		SearchField: andCriteria2Field,
		SearchType:  andCriteria2Type,
		SearchValue: andCriteria2Value,
	})

	updateReq := types.CreateOrUpdateAssetGroupRequest{
		GroupName:           updatedGroupName,
		GroupType:           groupType,
		GroupDescription:    updatedGroupDescription,
		MembershipPredicate: membershipPredicate,
	}

	updateSuccess, err := client.UpdateAssetGroup(ctx, groupID, updateReq)
	if err != nil {
		t.Fatalf("error creating asset group: %s", err.Error())
	}
	if !updateSuccess {
		t.Fatalf("asset group update unsuccessful: %s", err.Error())
	}

	// Check
	assert.Equal(t, true, updateSuccess)
	updatedAssetGroups, err := client.ListAssetGroups(ctx, listReq)
	if err != nil {
		t.Fatalf("error fetching updated asset group: %s", err.Error())
	}

	expectedUpdatedFilter := append(expectedFilter,
		types.AssetGroupFilter{
			PrettyName: "AND",
			RenderType: "connector",
		},
		types.AssetGroupFilter{
			PrettyName: andCriteria2Value,
			RenderType: "value",
		},
	)

	assert.NotNil(t, updatedAssetGroups)
	assert.Len(t, updatedAssetGroups, 1)
	updatedAssetGroup := updatedAssetGroups[0]
	assert.Equal(t, groupID, updatedAssetGroup.ID)
	assert.Equal(t, updatedGroupName, updatedAssetGroup.Name)
	assert.Equal(t, groupType, updatedAssetGroup.Type)
	assert.Equal(t, updatedGroupDescription, updatedAssetGroup.Description)
	assert.Equal(t, membershipPredicate, updatedAssetGroup.MembershipPredicate)
	assert.Equal(t, expectedUpdatedFilter, updatedAssetGroup.Filter)
}
