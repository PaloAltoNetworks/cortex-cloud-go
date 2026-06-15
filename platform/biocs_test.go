// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
	platformTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- helpers ---------------------------------------------------------------

// biocInsertRequest mirrors the body shape /bioc/insert expects:
// `{ "request_data": [<BIOC>, ...] }`.
type biocInsertRequest struct {
	RequestData []platformTypes.BIOC `json:"request_data"`
}

// biocFilterRequest mirrors /bioc/get and /bioc/delete:
// `{ "request_data": { "filters": [...], ... } }`.
type biocFilterRequest struct {
	RequestData platformTypes.ListBIOCsRequest `json:"request_data"`
}

// biocDeleteRequest mirrors the delete body.
type biocDeleteRequest struct {
	RequestData platformTypes.DeleteBIOCsRequest `json:"request_data"`
}

// xqlIndicator wraps a raw XQL string as the JSON value /bioc/insert expects
// when is_xql=true. The endpoint accepts the indicator field as either a
// string (XQL) or an object (filter AST); the SDK keeps it as a
// json.RawMessage so callers can pass either.
func xqlIndicator(query string) json.RawMessage {
	b, _ := json.Marshal(query)
	return json.RawMessage(b)
}

// structuredIndicator returns a minimal filter-AST indicator object.
func structuredIndicator(t *testing.T) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"runOnCGO":          true,
		"investigationType": "PROCESS_EXECUTION_EVENT",
		"investigation": map[string]any{
			"PROCESS_EXECUTION_EVENT": map[string]any{
				"filter": map[string]any{
					"AND": []map[string]any{
						{
							"SEARCH_FIELD": "action_process_username",
							"SEARCH_TYPE":  "EQ",
							"SEARCH_VALUE": "test",
							"EXTRA_FIELDS": []any{},
							"isExtended":   false,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	return json.RawMessage(b)
}

// --- InsertBIOCs -----------------------------------------------------------

func TestClient_InsertBIOCs(t *testing.T) {
	t.Run("create returns added_objects with new rule_id", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+InsertBIOCsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			var req biocInsertRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData, 1)
			assert.Equal(t, "test-bioc", req.RequestData[0].Name)
			assert.Equal(t, enums.BIOCTypeExecution, req.RequestData[0].Type)
			assert.Equal(t, enums.BIOCSeverityLow, req.RequestData[0].Severity)
			// rule_id is omitempty and absent on create.
			assert.Zero(t, req.RequestData[0].RuleID)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [{"id": 42, "status": "Created a new bioc with the ID: 42 successfully"}],
				"updated_objects": [],
				"errors": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{
			{
				Name:      "test-bioc",
				Type:      enums.BIOCTypeExecution,
				Severity:  enums.BIOCSeverityLow,
				Status:    enums.BIOCStatusEnabled,
				IsXQL:     true,
				Indicator: xqlIndicator("dataset = xdr_data | filter event_type = 1"),
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.AddedObjects, 1)
		assert.Equal(t, 42, resp.AddedObjects[0].ID)
		assert.Contains(t, resp.AddedObjects[0].Status, "Created")
		assert.Empty(t, resp.UpdatedObjects)
		assert.Empty(t, resp.Errors)
	})

	t.Run("update returns updated_objects when rule_id is supplied", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			var req biocInsertRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData, 1)
			assert.Equal(t, 42, req.RequestData[0].RuleID)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [],
				"updated_objects": [{"id": 42, "status": "Updated a bioc with the ID: 42 successfully"}],
				"errors": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{
			{
				RuleID:    42,
				Name:      "renamed-bioc",
				Type:      enums.BIOCTypeExecution,
				Severity:  enums.BIOCSeverityLow,
				Status:    enums.BIOCStatusEnabled,
				IsXQL:     true,
				Indicator: xqlIndicator("dataset = xdr_data | filter event_type = 1"),
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.UpdatedObjects, 1)
		assert.Equal(t, 42, resp.UpdatedObjects[0].ID)
		assert.Empty(t, resp.AddedObjects)
		assert.Empty(t, resp.Errors)
	})

	t.Run("HTTP 400 with success-shape body recovers errors[] as typed response", func(t *testing.T) {
		// Regression guard: per-record validation failures on /bioc/insert
		// return HTTP 400 with the success body shape. InsertBIOCs must
		// re-parse the body and surface resp.Errors with no Go error.
		// Verified against a live tenant.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"added_objects": [],
				"updated_objects": [],
				"errors": [{"index": 0, "status": "Failed to create bioc due to: Missing the fields: ['comment']"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{{
			Name: "bad", Type: enums.BIOCTypeExecution, Severity: enums.BIOCSeverityLow,
		}})
		require.NoError(t, err, "per-record errors on 400 should NOT surface as a Go error")
		require.Len(t, resp.Errors, 1)
		assert.Equal(t, 0, resp.Errors[0].Index)
		assert.Contains(t, resp.Errors[0].Status, "Missing the fields")
		assert.Empty(t, resp.AddedObjects)
		assert.Empty(t, resp.UpdatedObjects)
	})

	t.Run("HTTP 200 with errors are decoded as {index,status} objects", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [],
				"updated_objects": [],
				"errors": [{"index": 1, "status": "Failed to create bioc due to: Invalid type"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{{Name: "x"}})
		require.NoError(t, err)
		require.Len(t, resp.Errors, 1)
		assert.Equal(t, 1, resp.Errors[0].Index)
		assert.Contains(t, resp.Errors[0].Status, "Invalid type")
	})

	t.Run("transport-level errors still surface", func(t *testing.T) {
		// Non-recoverable HTTP errors (e.g. 401 with reply-wrapped body)
		// should not be silently swallowed by the 400-recovery path.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"reply":{"err_code":401,"err_msg":"Unauthorized","err_extra":""}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		_, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{{Name: "x"}})
		require.Error(t, err)
	})

	t.Run("success body is read at top level (no reply wrapper)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [{"id": 1, "status": "ok"}],
				"updated_objects": [],
				"errors": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{{
			Name: "x", Type: enums.BIOCTypeExecution, Severity: enums.BIOCSeverityLow,
		}})
		require.NoError(t, err)
		require.Len(t, resp.AddedObjects, 1)
		assert.Equal(t, 1, resp.AddedObjects[0].ID)
	})

	t.Run("polymorphic indicator: object payload passes through", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req biocInsertRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData, 1)
			// The Indicator field should round-trip as a JSON object,
			// not a JSON string. Decode as map and probe a field.
			var obj map[string]any
			require.NoError(t, json.Unmarshal(req.RequestData[0].Indicator, &obj))
			assert.Equal(t, "PROCESS_EXECUTION_EVENT", obj["investigationType"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"added_objects": [{"id": 7, "status": "Created"}], "updated_objects": [], "errors": []}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertBIOCs(context.Background(), []platformTypes.BIOC{{
			Name: "structured", Type: enums.BIOCTypeExecution, Severity: enums.BIOCSeverityLow,
			Status: enums.BIOCStatusEnabled, IsXQL: false, Indicator: structuredIndicator(t),
		}})
		require.NoError(t, err)
		require.Len(t, resp.AddedObjects, 1)
	})
}

// --- ListBIOCs -------------------------------------------------------------

func TestClient_ListBIOCs(t *testing.T) {
	t.Run("returns objects with read-only fields populated", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListBIOCsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req biocFilterRequest
			require.NoError(t, json.Unmarshal(body, &req))
			assert.True(t, req.RequestData.ExtendedView)
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "name", req.RequestData.Filters[0].Field)
			assert.Equal(t, "EQ", req.RequestData.Filters[0].Operator)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects_type": "bioc",
				"objects": [{
					"rule_id": 57,
					"name": "test",
					"type": "EXECUTION",
					"severity": "SEV_040_HIGH",
					"status": "ENABLED",
					"comment": "audit",
					"is_xql": true,
					"indicator": "dataset = xdr_data | filter event_type = 1",
					"mitre_tactic_id_and_name": ["TA0001 - Initial Access"],
					"mitre_technique_id_and_name": ["T1059 - Command and Scripting Interpreter"],
					"creation_time": 1781000000000,
					"modification_time": 1781000001000,
					"source": "Public API user (key #30)",
					"number_of_issues": 7
				}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListBIOCs(context.Background(), platformTypes.ListBIOCsRequest{
			ExtendedView: true,
			Filters: []platformTypes.BIOCFilter{
				{Field: "name", Operator: "EQ", Value: "test"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.ObjectsCount)
		assert.Equal(t, "bioc", resp.ObjectsType)
		require.Len(t, resp.Objects, 1)

		obj := resp.Objects[0]
		assert.Equal(t, 57, obj.RuleID)
		assert.Equal(t, "test", obj.Name)
		assert.Equal(t, enums.BIOCTypeExecution, obj.Type)
		assert.Equal(t, enums.BIOCSeverityHigh, obj.Severity)
		// Status is server-preserved: pre-existing tenant records may
		// come back uppercase. Callers normalize for comparison.
		assert.Equal(t, enums.BIOCStatus("ENABLED"), obj.Status)
		assert.True(t, obj.IsXQL)
		// Indicator decodes to a JSON string for XQL — keep the surrounding quotes.
		assert.Equal(t, `"dataset = xdr_data | filter event_type = 1"`, string(obj.Indicator))
		// Read-only fields:
		assert.Equal(t, int64(1781000000000), obj.CreationTime)
		assert.Equal(t, int64(1781000001000), obj.ModificationTime)
		assert.Equal(t, "Public API user (key #30)", obj.Source)
		assert.Equal(t, 7, obj.NumberOfIssues)
	})

	t.Run("empty result decodes cleanly", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": [], "objects_type": "bioc"}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListBIOCs(context.Background(), platformTypes.ListBIOCsRequest{})
		require.NoError(t, err)
		assert.Zero(t, resp.ObjectsCount)
		assert.Empty(t, resp.Objects)
	})
}

// --- DeleteBIOCs -----------------------------------------------------------

func TestClient_DeleteBIOCs(t *testing.T) {
	t.Run("returns deleted rule_ids", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+DeleteBIOCsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req biocDeleteRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "rule_id", req.RequestData.Filters[0].Field)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 1, "objects": [496]}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		ids, err := client.DeleteBIOCs(context.Background(), platformTypes.DeleteBIOCsRequest{
			Filters: []platformTypes.BIOCFilter{
				{Field: "rule_id", Operator: "EQ", Value: 496},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []int{496}, ids)
	})

	t.Run("idempotent: empty filter result is not an error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": []}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		ids, err := client.DeleteBIOCs(context.Background(), platformTypes.DeleteBIOCsRequest{
			Filters: []platformTypes.BIOCFilter{
				{Field: "rule_id", Operator: "EQ", Value: 999999},
			},
		})
		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}

// --- FindBIOCByID ----------------------------------------------------------

func TestClient_FindBIOCByID(t *testing.T) {
	t.Run("queries by rule_id and returns match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			var raw struct {
				RequestData struct {
					Filters []struct {
						Field    string `json:"field"`
						Operator string `json:"operator"`
						Value    any    `json:"value"`
					} `json:"filters"`
					ExtendedView bool `json:"extended_view"`
				} `json:"request_data"`
			}
			require.NoError(t, json.Unmarshal(body, &raw))
			require.Len(t, raw.RequestData.Filters, 1)
			assert.Equal(t, "rule_id", raw.RequestData.Filters[0].Field)
			assert.Equal(t, "EQ", raw.RequestData.Filters[0].Operator)
			// JSON numbers decode to float64 in an `any` field.
			assert.Equal(t, float64(192), raw.RequestData.Filters[0].Value)
			assert.True(t, raw.RequestData.ExtendedView)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects": [{"rule_id": 192, "name": "x", "type": "EXECUTION", "severity": "SEV_040_HIGH", "status": "enabled", "is_xql": true, "indicator": "dataset = xdr_data"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindBIOCByID(context.Background(), 192)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 192, got.RuleID)
		assert.Equal(t, enums.BIOCTypeExecution, got.Type)
		assert.Equal(t, enums.BIOCSeverityHigh, got.Severity)
	})

	t.Run("returns nil when no match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": []}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindBIOCByID(context.Background(), 999)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

// --- FindBIOCByName --------------------------------------------------------

func TestClient_FindBIOCByName(t *testing.T) {
	t.Run("returns first matching record (names not unique)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req biocFilterRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "name", req.RequestData.Filters[0].Field)
			assert.Equal(t, "EQ", req.RequestData.Filters[0].Operator)
			assert.Equal(t, "shared-name", req.RequestData.Filters[0].Value)
			assert.True(t, req.RequestData.ExtendedView)

			w.WriteHeader(http.StatusOK)
			// Two BIOCs share the name on this tenant; the helper
			// returns the first match.
			fmt.Fprint(w, `{
				"objects_count": 2,
				"objects": [
					{"rule_id": 9, "name": "shared-name", "type": "EXECUTION", "severity": "SEV_020_LOW"},
					{"rule_id": 11, "name": "shared-name", "type": "EXECUTION", "severity": "SEV_020_LOW"}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindBIOCByName(context.Background(), "shared-name")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 9, got.RuleID)
	})

	t.Run("returns nil when no match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": [], "objects_type": "bioc"}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindBIOCByName(context.Background(), "missing")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
