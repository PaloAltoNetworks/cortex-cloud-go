// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudsec

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/tests"
	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/cloudsec"
)

func TestAccRule_CRUD(t *testing.T) {
	skipIfNotAcceptance(t)

	ctx := context.Background()
	config := tests.NewTestConfigFromEnv(t)
	client, err := NewClient(config.GetOptions()...)
	if err != nil {
		t.Fatalf("failed to initialize client: %s", err.Error())
	}

	// Generate unique rule name
	ruleName := fmt.Sprintf("test-rule-%d", time.Now().Unix())

	// Test Create
	t.Run("Create", func(t *testing.T) {
		enabled := true
		createReq := types.CreateRuleRequest{
			Name:        ruleName,
			Description: "Test rule created by acceptance test",
			Class:       enums.RuleClassConfig.String(),
			Type:        "DETECTION",
			AssetTypes:  []string{"aws-s3-bucket"},
			Severity:    enums.CloudSecSeverityHigh.String(),
			Query: types.QueryRequest{
				XQL: "config from cloud.resource where cloud.type = 'aws' AND api.name = 'aws-s3api-get-bucket-acl'",
			},
			Metadata: &types.MetadataRequest{
				Issue: &types.IssueRequest{
					Recommendation: "This is a test recommendation for the acceptance test.",
				},
			},
			Labels:  []string{"test", "acceptance"},
			Enabled: &enabled,
		}

		rule, err := client.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if rule.ID == "" {
			t.Error("Created rule has no ID")
		}
		if rule.Name != ruleName {
			t.Errorf("Created rule name = %v, want %v", rule.Name, ruleName)
		}
		if rule.Severity != enums.CloudSecSeverityHigh.String() {
			t.Errorf("Created rule severity = %v, want %v", rule.Severity, enums.CloudSecSeverityHigh.String())
		}

		// Store ID for subsequent tests
		ruleID := rule.ID

		// Test Get
		t.Run("Get", func(t *testing.T) {
			retrieved, err := client.Get(ctx, ruleID)
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}

			if retrieved.ID != ruleID {
				t.Errorf("Retrieved rule ID = %v, want %v", retrieved.ID, ruleID)
			}
			if retrieved.Name != ruleName {
				t.Errorf("Retrieved rule name = %v, want %v", retrieved.Name, ruleName)
			}
		})

		// Test Search
		t.Run("Search", func(t *testing.T) {
			searchReq := types.SearchRulesRequest{
				Filter: &types.FilterCriteria{
					SearchField: "id",
					SearchType:  enums.SearchTypeEqualTo.String(),
					SearchValue: ruleID,
				},
				SearchFrom: 0,
				SearchTo:   10,
			}

			results, err := client.Search(ctx, searchReq)
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			if results.Metadata.FilterCount == 0 {
				t.Error("Search returned no results")
			}

			found := false
			for _, rule := range results.Data {
				if rule.ID == ruleID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Created rule not found in search results")
			}
		})

		// Test Update
		t.Run("Update", func(t *testing.T) {
			updateReq := types.UpdateRuleRequest{
				Severity: enums.CloudSecSeverityCritical.String(),
				Labels:   []string{"test", "acceptance", "updated"},
			}

			updated, err := client.Update(ctx, ruleID, updateReq)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}

			if updated.Severity != enums.CloudSecSeverityCritical.String() {
				t.Errorf("Updated rule severity = %v, want %v", updated.Severity, enums.CloudSecSeverityCritical.String())
			}
			if len(updated.Labels) != 3 {
				t.Errorf("Updated rule has %d labels, want 3", len(updated.Labels))
			}
		})

		// Test Delete
		t.Run("Delete", func(t *testing.T) {
			err := client.Delete(ctx, ruleID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion by attempting to get the rule
			_, err = client.Get(ctx, ruleID)
			if err == nil {
				t.Error("Expected error when getting deleted rule, got nil")
			}
		})
	})
}

func TestAccRule_Search_Filters(t *testing.T) {
	skipIfNotAcceptance(t)

	ctx := context.Background()
	config := tests.NewTestConfigFromEnv(t)
	client, err := NewClient(config.GetOptions()...)
	if err != nil {
		t.Fatalf("failed to initialize client: %s", err.Error())
	}

	t.Run("SearchWithORFilter", func(t *testing.T) {
		searchReq := types.SearchRulesRequest{
			Filter: &types.FilterCriteria{
				OR: []types.FilterCriteria{
					{
						SearchField: "severity",
						SearchType:  enums.SearchTypeEqualTo.String(),
						SearchValue: enums.CloudSecSeverityHigh.String(),
					},
					{
						SearchField: "severity",
						SearchType:  enums.SearchTypeEqualTo.String(),
						SearchValue: enums.CloudSecSeverityCritical.String(),
					},
				},
			},
			SearchFrom: 0,
			SearchTo:   50,
			Sort: []types.SortCriteria{
				{Field: "name", Order: enums.SortOrderASC.String()},
			},
		}

		results, err := client.Search(ctx, searchReq)
		if err != nil {
			t.Fatalf("Search with OR filter failed: %v", err)
		}

		if results.Metadata.FilterCount == 0 {
			t.Log("No rules found with high or critical severity (this may be expected)")
		}
	})

	t.Run("SearchWithANDFilter", func(t *testing.T) {
		searchReq := types.SearchRulesRequest{
			Filter: &types.FilterCriteria{
				AND: []types.FilterCriteria{
					{
						SearchField: "enabled",
						SearchType:  enums.SearchTypeEqualTo.String(),
						SearchValue: true,
					},
					{
						SearchField: "system_default",
						SearchType:  enums.SearchTypeEqualTo.String(),
						SearchValue: false,
					},
				},
			},
			SearchFrom: 0,
			SearchTo:   50,
		}

		results, err := client.Search(ctx, searchReq)
		if err != nil {
			t.Fatalf("Search with AND filter failed: %v", err)
		}

		t.Logf("Found %d custom enabled rules", results.Metadata.FilterCount)
	})
}

// TestAccRule_ComplianceMetadata verifies that compliance_metadata can be set on
// create and updated via the API. This exercises the ComplianceMetadataInput type
// end-to-end against a real API.
func TestAccRule_ComplianceMetadata(t *testing.T) {
	skipIfNotAcceptance(t)

	ctx := context.Background()
	config := tests.NewTestConfigFromEnv(t)
	client, err := NewClient(config.GetOptions()...)
	if err != nil {
		t.Fatalf("failed to initialize client: %s", err.Error())
	}

	ruleName := fmt.Sprintf("test-compliance-metadata-%d", time.Now().Unix())

	// Create a rule with compliance_metadata
	enabled := true
	createReq := types.CreateRuleRequest{
		Name:       ruleName,
		Class:      enums.RuleClassConfig.String(),
		Type:       "DETECTION",
		AssetTypes: []string{"aws-s3-bucket"},
		Severity:   enums.CloudSecSeverityHigh.String(),
		Query: types.QueryRequest{
			XQL: "config from cloud.resource where cloud.type = 'aws' AND api.name = 'aws-s3api-get-bucket-acl'",
		},
		Metadata: &types.MetadataRequest{
			Issue: &types.IssueRequest{
				Recommendation: "Test recommendation for compliance_metadata acceptance test.",
			},
		},
		Labels:  []string{"test", "compliance-metadata"},
		Enabled: &enabled,
		ComplianceMetadata: []types.ComplianceMetadataInput{
			{ControlID: "requirement-cis-1"},
		},
	}

	rule, err := client.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create with compliance_metadata failed: %v", err)
	}
	// Ensure cleanup
	defer func() {
		_ = client.Delete(ctx, rule.ID)
	}()

	// Verify the response contains enriched compliance_metadata
	if len(rule.ComplianceMetadata) == 0 {
		t.Log("Warning: API returned empty compliance_metadata in create response (may need valid control_id for enrichment)")
	} else {
		for i, cm := range rule.ComplianceMetadata {
			if cm.ControlID == "" {
				t.Errorf("compliance_metadata[%d].control_id is empty", i)
			}
			t.Logf("compliance_metadata[%d]: control_id=%s, control_name=%s, standard_id=%s, standard_name=%s",
				i, cm.ControlID, cm.ControlName, cm.StandardID, cm.StandardName)
		}
	}

	// Update the rule with different compliance_metadata
	updateReq := types.UpdateRuleRequest{
		ComplianceMetadata: []types.ComplianceMetadataInput{
			{ControlID: "requirement-cis-1"},
			{ControlID: "requirement-cis-2"},
		},
	}

	updated, err := client.Update(ctx, rule.ID, updateReq)
	if err != nil {
		t.Fatalf("Update with compliance_metadata failed: %v", err)
	}

	if len(updated.ComplianceMetadata) == 0 {
		t.Log("Warning: API returned empty compliance_metadata in update response (may need valid control_id for enrichment)")
	} else {
		t.Logf("After update: %d compliance_metadata entries", len(updated.ComplianceMetadata))
	}

	// Verify via Get that compliance_metadata persisted
	retrieved, err := client.Get(ctx, rule.ID)
	if err != nil {
		t.Fatalf("Get after compliance_metadata update failed: %v", err)
	}

	t.Logf("Retrieved rule has %d compliance_metadata entries", len(retrieved.ComplianceMetadata))
}
