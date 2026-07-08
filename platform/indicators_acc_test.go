// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package platform

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
	platformTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccIndicatorLifecycle exercises the full create→read→update→delete
// path against a live tenant. Configure via the same env vars used by the
// asset-group acceptance test: TEST_CORTEX_API_URL, TEST_CORTEX_API_KEY,
// TEST_CORTEX_API_KEY_ID. Run with:
//
//	go test -tags=acceptance -run TestAccIndicator -v ./platform/...
func TestAccIndicatorLifecycle(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	indicatorName := fmt.Sprintf("go-sdk-acctest-%s.example.test", timestamp)
	original := platformTypes.Indicator{
		Indicator:                indicatorName,
		Type:                     enums.IndicatorTypeDomainName,
		Severity:                 enums.IndicatorSeverityLow,
		ExpirationDate:           -1,
		DefaultExpirationEnabled: true,
		Comment:                  fmt.Sprintf("go-sdk acceptance test %s", timestamp),
		Reputation:               enums.IndicatorReputationUnknown,
		Reliability:              enums.IndicatorReliabilityF,
	}

	// --- Create ---
	insertResp, err := client.InsertIndicators(ctx, []platformTypes.Indicator{original})
	require.NoError(t, err, "InsertIndicators (create) failed")
	require.Empty(t, insertResp.Errors, "unexpected per-record errors: %v", insertResp.Errors)
	require.Len(t, insertResp.AddedObjects, 1, "expected one added record")
	require.Empty(t, insertResp.UpdatedObjects, "create path should not populate updated_objects")
	ruleID := insertResp.AddedObjects[0].ID
	assert.Positive(t, ruleID)
	assert.Contains(t, insertResp.AddedObjects[0].Status, "Created")

	// Defer cleanup as early as possible so a panic inside this test
	// doesn't leak the record.
	defer func() {
		ids, delErr := client.DeleteIndicators(ctx, platformTypes.DeleteIndicatorsRequest{
			Filters: []platformTypes.IndicatorFilter{
				{Field: "indicator", Operator: "EQ", Value: indicatorName},
			},
		})
		if delErr != nil {
			t.Logf("cleanup delete failed: %s", delErr.Error())
		}
		// Tolerate already-gone — the explicit delete near the end of
		// the test may have run first.
		t.Logf("cleanup removed rule_ids: %v", ids)
	}()

	// --- Read by name ---
	byName, err := client.FindIndicatorByName(ctx, indicatorName)
	require.NoError(t, err, "FindIndicatorByName failed")
	require.NotNil(t, byName, "indicator should exist after insert")
	assert.Equal(t, ruleID, byName.RuleID)
	assert.Equal(t, indicatorName, byName.Indicator)
	assert.Equal(t, enums.IndicatorTypeDomainName, byName.Type)
	assert.Equal(t, enums.IndicatorSeverityLow, byName.Severity)
	assert.Equal(t, int64(-1), byName.ExpirationDate)
	assert.True(t, byName.DefaultExpirationEnabled)
	assert.Equal(t, original.Comment, byName.Comment)
	assert.Equal(t, enums.IndicatorReputationUnknown, byName.Reputation)
	assert.Equal(t, enums.IndicatorReliabilityF, byName.Reliability)
	// Read-only fields (S3 regression guard): the live API populates
	// these when extended_view=true; FindIndicatorByName already requests
	// extended_view.
	assert.Positive(t, byName.CreationTime, "creation_time should be a unix-epoch-ms")
	assert.Positive(t, byName.ModificationTime, "modification_time should be a unix-epoch-ms")
	assert.NotEmpty(t, byName.Status, "status (e.g. ENABLED) should be set")
	assert.NotEmpty(t, byName.Source, "source (e.g. Public API user (key #N)) should be set")
	assert.GreaterOrEqual(t, byName.NumberOfIssues, 0)

	// --- Read by ID ---
	// rule_id filter is undocumented in the OpenAPI field enum but works
	// against the live API on EQ. This sub-assertion is the canonical
	// regression guard for B2.
	byID, err := client.FindIndicatorByID(ctx, ruleID)
	require.NoError(t, err, "FindIndicatorByID failed")
	require.NotNil(t, byID, "rule_id lookup should succeed")
	assert.Equal(t, ruleID, byID.RuleID)
	assert.Equal(t, indicatorName, byID.Indicator)

	// --- Update via upsert (rule_id-keyed overwrite) ---
	updated := original
	updated.RuleID = ruleID
	updated.Severity = enums.IndicatorSeverityCritical
	updated.Reputation = enums.IndicatorReputationBad
	updated.Reliability = enums.IndicatorReliabilityA
	updated.Comment = original.Comment + " (updated)"

	updateResp, err := client.InsertIndicators(ctx, []platformTypes.Indicator{updated})
	require.NoError(t, err, "InsertIndicators (update) failed")
	require.Empty(t, updateResp.Errors, "unexpected per-record errors on update: %v", updateResp.Errors)
	require.Empty(t, updateResp.AddedObjects, "update path should not populate added_objects")
	require.Len(t, updateResp.UpdatedObjects, 1, "expected one updated record")
	assert.Equal(t, ruleID, updateResp.UpdatedObjects[0].ID)
	assert.Contains(t, updateResp.UpdatedObjects[0].Status, "Updated")

	// Confirm new content lands on read-back.
	afterUpdate, err := client.FindIndicatorByName(ctx, indicatorName)
	require.NoError(t, err)
	require.NotNil(t, afterUpdate)
	assert.Equal(t, ruleID, afterUpdate.RuleID, "rule_id should be stable across in-place updates")
	assert.Equal(t, enums.IndicatorSeverityCritical, afterUpdate.Severity)
	assert.Equal(t, enums.IndicatorReputationBad, afterUpdate.Reputation)
	assert.Equal(t, enums.IndicatorReliabilityA, afterUpdate.Reliability)
	assert.Equal(t, updated.Comment, afterUpdate.Comment)
	assert.GreaterOrEqual(t, afterUpdate.ModificationTime, byName.ModificationTime,
		"modification_time should be non-decreasing after an update")

	// --- Delete by name ---
	deleted, err := client.DeleteIndicators(ctx, platformTypes.DeleteIndicatorsRequest{
		Filters: []platformTypes.IndicatorFilter{
			{Field: "indicator", Operator: "EQ", Value: indicatorName},
		},
	})
	require.NoError(t, err)
	require.Len(t, deleted, 1, "expected exactly one record removed")
	assert.Equal(t, ruleID, deleted[0])

	// Confirm gone.
	gone, err := client.FindIndicatorByName(ctx, indicatorName)
	require.NoError(t, err)
	assert.Nil(t, gone, "indicator should be absent after delete")
}

// TestAccIndicatorIdempotentDelete pins C1: /indicators/delete returns
// `{objects_count: 0, objects: []}` (no error) when the filter matches
// nothing.
func TestAccIndicatorIdempotentDelete(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	missing := fmt.Sprintf("go-sdk-acctest-missing-%s.example.test", timestamp)
	ids, err := client.DeleteIndicators(ctx, platformTypes.DeleteIndicatorsRequest{
		Filters: []platformTypes.IndicatorFilter{
			{Field: "indicator", Operator: "EQ", Value: missing},
		},
	})
	require.NoError(t, err, "deleting a non-existent indicator should not error")
	assert.Empty(t, ids, "expected zero rule_ids returned for missing record")
}

// TestAccIndicatorURLType verifies B1 — `type=URL` is accepted by the live
// API even though it's omitted from the OpenAPI insert enum. SEV_050_CRITICAL
// is verified at the same time.
func TestAccIndicatorURLType(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	indicatorURL := fmt.Sprintf("https://go-sdk-acctest-%s.example.test/probe", timestamp)
	resp, err := client.InsertIndicators(ctx, []platformTypes.Indicator{{
		Indicator:                indicatorURL,
		Type:                     enums.IndicatorTypeURL,
		Severity:                 enums.IndicatorSeverityCritical,
		ExpirationDate:           -1,
		DefaultExpirationEnabled: true,
		Reputation:               enums.IndicatorReputationBad,
		Reliability:              enums.IndicatorReliabilityA,
		Comment:                  "URL + SEV_050_CRITICAL acceptance probe",
	}})
	require.NoError(t, err)
	require.Empty(t, resp.Errors, "URL type should be accepted: %v", resp.Errors)
	require.Len(t, resp.AddedObjects, 1)
	ruleID := resp.AddedObjects[0].ID

	defer func() {
		_, _ = client.DeleteIndicators(ctx, platformTypes.DeleteIndicatorsRequest{
			Filters: []platformTypes.IndicatorFilter{
				{Field: "indicator", Operator: "EQ", Value: indicatorURL},
			},
		})
	}()

	got, err := client.FindIndicatorByID(ctx, ruleID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, enums.IndicatorTypeURL, got.Type)
	assert.Equal(t, enums.IndicatorSeverityCritical, got.Severity)
}

// TestAccIndicatorBooleanFilter verifies L2 — the
// `default_expiration_enabled` filter requires a JSON boolean. Sending the
// value as a JSON bool (the SDK shape) succeeds; this test pins the wire
// type contract end-to-end. The negative case (string would 500) lives in
// the unit tests where it can be exercised without touching the tenant.
func TestAccIndicatorBooleanFilter(t *testing.T) {
	client := setupAcceptanceTest(t)
	ctx := context.Background()

	resp, err := client.ListIndicators(ctx, platformTypes.ListIndicatorsRequest{
		ExtendedView: true,
		Filters: []platformTypes.IndicatorFilter{
			{Field: "default_expiration_enabled", Operator: "EQ", Value: true},
		},
		SearchTo: 1,
	})
	require.NoError(t, err, "boolean filter sent as JSON bool should not error")
	// We can't assert a specific count (depends on tenant state), only
	// that the call decodes cleanly into the SDK type.
	assert.GreaterOrEqual(t, resp.ObjectsCount, 0)
}
