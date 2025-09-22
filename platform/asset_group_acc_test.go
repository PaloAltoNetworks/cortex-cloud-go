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

	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func setupAcceptanceTest(t *testing.T) *Client {
	apiUrl := os.Getenv("CORTEX_API_URL_TEST")
	apiKey := os.Getenv("CORTEX_API_KEY_TEST")
	apiKeyIDStr := os.Getenv("CORTEX_API_KEY_ID_TEST")

	apiKeyID, err := strconv.Atoi(apiKeyIDStr)
	if err != nil {
		t.Fatalf("failed to convert API key ID \"%s\" to int: %s", apiKeyIDStr, err.Error())
	}
	
	config := &api.Config{
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

//func TestAccStaticAssetGroupLifecycle(t *testing.T) {
//}

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
	andCriteria1Value := "test" // TODO: create an asset beforehand and set this to it
	membershipPredicate := types.CriteriaFilter{
		And: []types.Criteria{
			{
				SearchField: andCriteria1Field,
				SearchType: andCriteria1Type,
				SearchValue: andCriteria1Value,
			},
		},
	}

	createReq := CreateAssetGroupRequest{
		GroupName: groupName,
		GroupType: groupType,
		GroupDescription: groupDescription,
		MembershipPredicate: membershipPredicate,
	}

	createResp, err := client.CreateAssetGroup(ctx, createReq)
	if err != nil {
		t.Fatalf("error creating asset group: %s", err.Error())
	}
	assert.NotNil(t, createResp)
	assert.Equal(t, true, createResp.Success)
	assert.Positive(t, createResp.AssetGroupID)
	groupID := createResp.AssetGroupID

	// Read
	listReq := ListAssetGroupsRequest{
		Filters: types.CriteriaFilter{
			And: []types.Criteria{
				{
					SearchField: "XDM.ASSET_GROUP.NAME",
					SearchType: "CONTAINS",
					SearchValue: groupName,
				},
			},
		},
	}

	// Check
	listResp, err := client.ListAssetGroups(ctx, listReq)	
	if err != nil {
		t.Fatalf("error fetching asset group: %s", err.Error())
	}
	expectedFilter := []AssetGroupFilter{
		{
			PrettyName: "name",
			DataType: "TEXT",
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

	assert.NotNil(t, listResp)
	assert.Len(t, listResp.Data, 1)
	listRespData := listResp.Data[0]
	assert.Equal(t, groupID, listRespData.ID)
	assert.Equal(t, groupName, listRespData.Name)
	assert.Equal(t, groupType, listRespData.Type)
	assert.Equal(t, groupDescription, listRespData.Description)
	assert.Equal(t, membershipPredicate, listRespData.MembershipPredicate)
	assert.Equal(t, expectedFilter, listRespData.Filter)

	// Update
	// TODO:

	// Delete
	deleteResp, err := client.DeleteAssetGroup(ctx, groupID)
	if err != nil {
		t.Fatalf("error deleting asset group: %s", err.Error())
	}
	assert.NotNil(t, deleteResp)
	assert.Equal(t, true, deleteResp.Success)
}
