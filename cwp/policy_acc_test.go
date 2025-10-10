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

	"github.com/PaloAltoNetworks/cortex-cloud-go/types/cwp"
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

	client, err := NewClient(
		WithCortexAPIURL(apiUrl),
		WithCortexAPIKey(apiKey),
		WithCortexAPIKeyID(apiKeyID),
		WithLogLevel("debug"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	return client
}

func TestAccMisconfigurationPolicyLifecycle(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	currentTime := time.Now()
	timestamp := strconv.FormatInt(currentTime.Unix(), 10)
	name := fmt.Sprintf("go-sdk-acc-test-misconfig-%s", timestamp)
	createReq := types.CreatePolicyRequest{
		PolicyAction: "ISSUE",
		// Of the following two fields, which is required (or are both required)?
		//AssetGroups: []string{ "agent_compliance" },
		AssetGroupIDs: []int{19},
		//AssetScope: []byte{'e', 'y', 'J', 'B', 'T', 'k', 'Q', 'i', 'O', 'l', 't', '7', 'I', 'l', 'N', 'F', 'Q', 'V', 'J', 'D', 'S', 'F', '9', 'G', 'S', 'U', 'V', 'M', 'R', 'C', 'I', '6', 'I', 'n', 'h', 'k', 'b', 'S', '5', 'h', 'c', '3', 'N', 'l', 'd', 'C', '5', 'u', 'Y', 'W', '1', 'l', 'I', 'i', 'w', 'i', 'U', '0', 'V', 'B', 'U', 'k', 'N', 'I', 'X', '1', 'R', 'Z', 'U', 'E', 'U', 'i', 'O', 'i', 'J', 'D', 'T', '0', '5', 'U', 'Q', 'U', 'l', 'O', 'U', 'y', 'I', 's', 'I', 'l', 'N', 'F', 'Q', 'V', 'J', 'D', 'S', 'F', '9', 'W', 'Q', 'U', 'x', 'V', 'R', 'S', 'I', '6', 'I', 'm', 'N', 'v', 'b', 'X', 'B', 's', 'a', 'W', 'F', 'u', 'Y', '2', 'U', 't', 'Y', 'X', 'V', '0', 'b', '2', '1', 'h', 'd', 'G', 'l', 'v', 'b', 'i', '1', 'n', 'Y', '3', 'A', 'i', 'f', 'V', '1', '9', },
		// Why are the following values byte arrays instead of just strings?
		//Condition: []byte{},
		//Exception: []byte{},
		//CreatedAt: currentTime.String(),
		Description: "Go SDK Test",
		//Disabled: false,
		//EvaluationModes: []string{},
		EvaluationStage: "RUNTIME",
		//ModifiedAt: currentTime.String(),
		Name:                name,
		RemediationGuidance: "Test guidance",
		//Revision: 0,
		RulesIDs: []string{
			"00000000-0000-0000-0000-000000300419",
		},
		PolicySeverity: "CRITICAL",
		Type:           "COMPLIANCE",
	}

	createResp, err := client.CreatePolicy(ctx, createReq)
	if err != nil {
		t.Fatalf("error creating policy: %s", err.Error())
	}

	// Check
	assert.NotNil(t, createResp)
	//	assert.Equal(t, true, createResp.Success)
	//	assert.Positive(t, createResp.AssetGroupID)
	//	groupID := createResp.AssetGroupID

	//// Defer Delete check
	//testDeleteAssetGroup := func(t *testing.T, ctx context.Context, groupID int) {
	//	// Delete
	//	deleteResp, err := client.DeleteAssetGroup(ctx, groupID)

	//	// Check
	//	if err != nil {
	//		t.Fatalf("error deleting asset group: %s", err.Error())
	//	}
	//	assert.NotNil(t, deleteResp)
	//	assert.Equal(t, true, deleteResp.Success)
	//}
	//defer testDeleteAssetGroup(t, ctx, groupID)

	//// Read
	//listReq := types.ListAssetGroupsRequest{
	//	Filters: search.CriteriaFilter{
	//		And: []search.Criteria{
	//			{
	//				SearchField: "XDM.ASSET_GROUP.NAME",
	//				SearchType: "CONTAINS",
	//				SearchValue: groupName,
	//			},
	//		},
	//	},
	//}
	//listResp, err := client.ListAssetGroups(ctx, listReq)
	//if err != nil {
	//	t.Fatalf("error fetching asset group: %s", err.Error())
	//}

	//// Check
	//expectedFilter := []types.AssetGroupFilter{
	//	{
	//		PrettyName: "name",
	//		DataType: "TEXT",
	//		RenderType: "attribute",
	//	},
	//	{
	//		PrettyName: "doesn't contain",
	//		RenderType: "operator",
	//	},
	//	{
	//		PrettyName: andCriteria1Value,
	//		RenderType: "value",
	//	},
	//}

	//assert.NotNil(t, listResp)
	//assert.Len(t, listResp.Data, 1)
	//listRespData := listResp.Data[0]
	//assert.Equal(t, groupID, listRespData.ID)
	//assert.Equal(t, groupName, listRespData.Name)
	//assert.Equal(t, groupType, listRespData.Type)
	//assert.Equal(t, groupDescription, listRespData.Description)
	//assert.Equal(t, membershipPredicate, listRespData.MembershipPredicate)
	//assert.Equal(t, expectedFilter, listRespData.Filter)

	//// Update
	//updatedGroupName := fmt.Sprintf("%s-updated", groupName)
	//updatedGroupDescription := fmt.Sprintf("%s (Updated)", groupDescription)
	//andCriteria2Field := "xdm.asset.name"
	//andCriteria2Type := "NCONTAINS"
	//andCriteria2Value := "test"
	//membershipPredicate.And = append(membershipPredicate.And, search.Criteria{
	//	SearchField: andCriteria2Field,
	//	SearchType: andCriteria2Type,
	//	SearchValue: andCriteria2Value,
	//})

	//updateReq := types.CreateOrUpdateAssetGroupRequest{
	//	GroupName: updatedGroupName,
	//	GroupType: groupType,
	//	GroupDescription: updatedGroupDescription,
	//	MembershipPredicate: membershipPredicate,
	//}
	//
	//updateResp, err := client.UpdateAssetGroup(ctx, groupID, updateReq)
	//if err != nil {
	//	t.Fatalf("error creating asset group: %s", err.Error())
	//}

	//// Check
	//assert.NotNil(t, updateResp)
	//assert.Equal(t, true, updateResp.Success)
	//listResp, err = client.ListAssetGroups(ctx, listReq)
	//if err != nil {
	//	t.Fatalf("error fetching updated asset group: %s", err.Error())
	//}

	//expectedUpdatedFilter := append(expectedFilter,
	//	types.AssetGroupFilter{
	//		PrettyName: "AND",
	//		RenderType: "connector",
	//	},
	//	types.AssetGroupFilter{
	//		PrettyName: andCriteria2Value,
	//		RenderType: "value",
	//	},
	//)

	//assert.NotNil(t, listResp)
	//assert.Len(t, listResp.Data, 1)
	//listRespData = listResp.Data[0]
	//assert.Equal(t, groupID, listRespData.ID)
	//assert.Equal(t, updatedGroupName, listRespData.Name)
	//assert.Equal(t, groupType, listRespData.Type)
	//assert.Equal(t, updatedGroupDescription, listRespData.Description)
	//assert.Equal(t, membershipPredicate, listRespData.MembershipPredicate)
	//assert.Equal(t, expectedUpdatedFilter, listRespData.Filter)
}
