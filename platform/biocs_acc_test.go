// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
	platformTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccBIOCLifecycle exercises the full create→read→update→delete path
// against a live tenant. Configure via the same env vars used by the
// asset-group and indicator acceptance tests: TEST_CORTEX_API_URL,
// TEST_CORTEX_API_KEY, TEST_CORTEX_API_KEY_ID. Run with:
//
//	go test -tags=acceptance -run TestAccBIOC -v ./platform/...
func TestAccBIOCLifecycle(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	biocName := fmt.Sprintf("go-sdk-acctest-bioc-%s", timestamp)
	xqlBody, err := json.Marshal(fmt.Sprintf(
		"dataset = xdr_data | filter event_type = 1 and actor_process_image_name = \"acctest-%s.exe\"",
		timestamp,
	))
	require.NoError(t, err)

	original := platformTypes.BIOC{
		Name:                    biocName,
		Type:                    enums.BIOCTypeExecution,
		Severity:                enums.BIOCSeverityLow,
		Status:                  enums.BIOCStatusEnabled,
		Comment:                 fmt.Sprintf("go-sdk acceptance test %s", timestamp),
		IsXQL:                   true,
		Indicator:               json.RawMessage(xqlBody),
		MitreTacticIDAndName:    []string{},
		MitreTechniqueIDAndName: []string{},
	}

	// --- Create ---
	insertResp, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{original})
	require.NoError(t, err, "InsertBIOCs (create) failed")
	require.Empty(t, insertResp.Errors, "unexpected per-record errors: %v", insertResp.Errors)
	require.Len(t, insertResp.AddedObjects, 1, "expected one added record")
	require.Empty(t, insertResp.UpdatedObjects, "create path should not populate updated_objects")
	ruleID := insertResp.AddedObjects[0].ID
	assert.Positive(t, ruleID)
	assert.Contains(t, insertResp.AddedObjects[0].Status, "Created")

	// Defer cleanup as early as possible.
	defer func() {
		ids, delErr := client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
			Filters: []platformTypes.BIOCFilter{
				{Field: "rule_id", Operator: "EQ", Value: ruleID},
			},
		})
		if delErr != nil {
			t.Logf("cleanup delete failed: %s", delErr.Error())
		}
		t.Logf("cleanup removed rule_ids: %v", ids)
	}()

	// --- Read by ID ---
	// rule_id filter on /bioc/get is undocumented in the OpenAPI field
	// enum but accepted by the live API on EQ. This sub-assertion is the
	// canonical regression guard.
	byID, err := client.FindBIOCByID(ctx, ruleID)
	require.NoError(t, err, "FindBIOCByID failed")
	require.NotNil(t, byID, "rule_id lookup should succeed")
	assert.Equal(t, ruleID, byID.RuleID)
	assert.Equal(t, biocName, byID.Name)
	assert.Equal(t, enums.BIOCTypeExecution, byID.Type)
	assert.Equal(t, enums.BIOCSeverityLow, byID.Severity)
	assert.True(t, byID.IsXQL)
	assert.Positive(t, byID.CreationTime, "creation_time should be set")
	assert.Positive(t, byID.ModificationTime, "modification_time should be set")
	assert.NotEmpty(t, byID.Source, "source should be set")

	// --- Read by name ---
	byName, err := client.FindBIOCByName(ctx, biocName)
	require.NoError(t, err, "FindBIOCByName failed")
	require.NotNil(t, byName, "BIOC should exist after insert")
	assert.Equal(t, ruleID, byName.RuleID)

	// --- Update via upsert (rule_id-keyed overwrite) ---
	updated := original
	updated.RuleID = ruleID
	updated.Name = biocName + "-renamed" // rename in place: should keep rule_id
	updated.Severity = enums.BIOCSeverityHigh
	updated.Comment = original.Comment + " (updated)"

	updateResp, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{updated})
	require.NoError(t, err, "InsertBIOCs (update) failed")
	require.Empty(t, updateResp.Errors, "unexpected per-record errors on update: %v", updateResp.Errors)
	require.Empty(t, updateResp.AddedObjects, "update path should not populate added_objects")
	require.Len(t, updateResp.UpdatedObjects, 1, "expected one updated record")
	assert.Equal(t, ruleID, updateResp.UpdatedObjects[0].ID)
	assert.Contains(t, updateResp.UpdatedObjects[0].Status, "Updated")

	// Confirm new content lands on read-back.
	afterUpdate, err := client.FindBIOCByID(ctx, ruleID)
	require.NoError(t, err)
	require.NotNil(t, afterUpdate)
	assert.Equal(t, ruleID, afterUpdate.RuleID, "rule_id should be stable across in-place updates")
	assert.Equal(t, updated.Name, afterUpdate.Name, "name should be mutable via rule_id-keyed upsert")
	assert.Equal(t, enums.BIOCSeverityHigh, afterUpdate.Severity)
	assert.Equal(t, updated.Comment, afterUpdate.Comment)
	assert.GreaterOrEqual(t, afterUpdate.ModificationTime, byID.ModificationTime,
		"modification_time should be non-decreasing after an update")

	// --- Delete by rule_id ---
	deleted, err := client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
		Filters: []platformTypes.BIOCFilter{
			{Field: "rule_id", Operator: "EQ", Value: ruleID},
		},
	})
	require.NoError(t, err)
	require.Len(t, deleted, 1, "expected exactly one record removed")
	assert.Equal(t, ruleID, deleted[0])

	// Confirm gone.
	gone, err := client.FindBIOCByID(ctx, ruleID)
	require.NoError(t, err)
	assert.Nil(t, gone, "BIOC should be absent after delete")
}

// TestAccBIOCIdempotentDelete pins the empty-result-no-error contract:
// /bioc/delete returns `{objects_count: 0, objects: []}` (no error) when
// the filter matches nothing.
func TestAccBIOCIdempotentDelete(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()

	ids, err := client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
		Filters: []platformTypes.BIOCFilter{
			{Field: "rule_id", Operator: "EQ", Value: 999999999},
		},
	})
	require.NoError(t, err, "deleting a non-existent BIOC should not error")
	assert.Empty(t, ids, "expected zero rule_ids returned for missing record")
}

// TestAccBIOCStructuredIndicator verifies the non-XQL form: the indicator
// field is a JSON object describing a filter AST. The structured form is
// what the Cortex UI emits when users build a BIOC via the visual editor.
func TestAccBIOCStructuredIndicator(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	biocName := fmt.Sprintf("go-sdk-acctest-bioc-structured-%s", timestamp)
	indicator, err := json.Marshal(map[string]any{
		"runOnCGO":          true,
		"investigationType": "PROCESS_EXECUTION_EVENT",
		"investigation": map[string]any{
			"PROCESS_EXECUTION_EVENT": map[string]any{
				"filter": map[string]any{
					"AND": []map[string]any{
						{
							"SEARCH_FIELD": "action_process_username",
							"SEARCH_TYPE":  "EQ",
							"SEARCH_VALUE": fmt.Sprintf("acctest-%s", timestamp),
							"EXTRA_FIELDS": []any{},
							"isExtended":   false,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	resp, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{{
		Name:                    biocName,
		Type:                    enums.BIOCTypeExecution,
		Severity:                enums.BIOCSeverityInfo,
		Status:                  enums.BIOCStatusDisabled,
		Comment:                 "structured-indicator acceptance probe — safe to delete",
		IsXQL:                   false,
		Indicator:               json.RawMessage(indicator),
		MitreTacticIDAndName:    []string{},
		MitreTechniqueIDAndName: []string{},
	}})
	require.NoError(t, err)
	require.Empty(t, resp.Errors, "structured indicator should be accepted: %v", resp.Errors)
	require.Len(t, resp.AddedObjects, 1)
	ruleID := resp.AddedObjects[0].ID

	defer func() {
		_, _ = client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
			Filters: []platformTypes.BIOCFilter{
				{Field: "rule_id", Operator: "EQ", Value: ruleID},
			},
		})
	}()

	got, err := client.FindBIOCByID(ctx, ruleID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.False(t, got.IsXQL)

	// The indicator round-trips as a JSON object — decode and verify a
	// nested field survived.
	var obj map[string]any
	require.NoError(t, json.Unmarshal(got.Indicator, &obj))
	assert.Equal(t, "PROCESS_EXECUTION_EVENT", obj["investigationType"])
}

// TestAccBIOCSeverityCritical verifies that SEV_050_CRITICAL is accepted by
// /bioc/insert and round-trips via /bioc/get even though the OpenAPI
// severity enum tops at SEV_040_HIGH. Mirrors TestAccIndicatorURLType for
// IndicatorSeverityCritical.
func TestAccBIOCSeverityCritical(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	biocName := fmt.Sprintf("go-sdk-acctest-bioc-critical-%s", timestamp)
	xqlBody, err := json.Marshal(fmt.Sprintf(
		"dataset = xdr_data | filter event_type = 1 and actor_process_image_name = \"acctest-critical-%s.exe\"",
		timestamp,
	))
	require.NoError(t, err)

	resp, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{{
		Name:                    biocName,
		Type:                    enums.BIOCTypeExecution,
		Severity:                enums.BIOCSeverityCritical,
		Status:                  enums.BIOCStatusDisabled,
		Comment:                 "SEV_050_CRITICAL acceptance probe — safe to delete",
		IsXQL:                   true,
		Indicator:               json.RawMessage(xqlBody),
		MitreTacticIDAndName:    []string{},
		MitreTechniqueIDAndName: []string{},
	}})
	require.NoError(t, err)
	require.Empty(t, resp.Errors, "SEV_050_CRITICAL should be accepted: %v", resp.Errors)
	require.Len(t, resp.AddedObjects, 1)
	ruleID := resp.AddedObjects[0].ID

	defer func() {
		_, _ = client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
			Filters: []platformTypes.BIOCFilter{
				{Field: "rule_id", Operator: "EQ", Value: ruleID},
			},
		})
	}()

	got, err := client.FindBIOCByID(ctx, ruleID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, enums.BIOCSeverityCritical, got.Severity, "SEV_050_CRITICAL should round-trip via /bioc/get")
}

// TestAccBIOCNonUniqueNames pins the verified-against-live-tenant fact that
// BIOC names are NOT unique per tenant: two BIOCs with the same name
// produce two distinct rule_ids. This is the foundational invariant the
// terraform-provider-cortexcloud BIOC resource relies on for its
// rule_id-as-identity design.
func TestAccBIOCNonUniqueNames(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	biocName := fmt.Sprintf("go-sdk-acctest-bioc-shared-%s", timestamp)
	xqlA, _ := json.Marshal(fmt.Sprintf("dataset = xdr_data | filter event_type = 1 and actor_process_image_name = \"acctest-a-%s.exe\"", timestamp))
	xqlB, _ := json.Marshal(fmt.Sprintf("dataset = xdr_data | filter event_type = 1 and actor_process_image_name = \"acctest-b-%s.exe\"", timestamp))

	mkBIOC := func(body []byte) platformTypes.BIOC {
		return platformTypes.BIOC{
			Name:                    biocName,
			Type:                    enums.BIOCTypeExecution,
			Severity:                enums.BIOCSeverityInfo,
			Status:                  enums.BIOCStatusDisabled,
			Comment:                 "non-unique-name acceptance probe — safe to delete",
			IsXQL:                   true,
			Indicator:               json.RawMessage(body),
			MitreTacticIDAndName:    []string{},
			MitreTechniqueIDAndName: []string{},
		}
	}

	respA, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{mkBIOC(xqlA)})
	require.NoError(t, err)
	require.Len(t, respA.AddedObjects, 1)
	idA := respA.AddedObjects[0].ID

	respB, err := client.InsertBIOCs(ctx, []platformTypes.BIOC{mkBIOC(xqlB)})
	require.NoError(t, err)
	require.Len(t, respB.AddedObjects, 1)
	idB := respB.AddedObjects[0].ID

	defer func() {
		for _, id := range []int{idA, idB} {
			_, _ = client.DeleteBIOCs(ctx, platformTypes.DeleteBIOCsRequest{
				Filters: []platformTypes.BIOCFilter{
					{Field: "rule_id", Operator: "EQ", Value: id},
				},
			})
		}
	}()

	assert.NotEqual(t, idA, idB, "two BIOCs with the same name must get distinct rule_ids")
}
