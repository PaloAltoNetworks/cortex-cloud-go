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

// indicatorInsertRequest mirrors the body shape /indicators/insert expects:
// `{ "request_data": [<Indicator>, ...] }`.
type indicatorInsertRequest struct {
	RequestData []platformTypes.Indicator `json:"request_data"`
}

// indicatorFilterRequest mirrors /indicators/get and /indicators/delete:
// `{ "request_data": { "filters": [...], ... } }`.
type indicatorFilterRequest struct {
	RequestData platformTypes.ListIndicatorsRequest `json:"request_data"`
}

// indicatorDeleteRequest mirrors the delete body — Filters is the only
// member that gets populated by the SDK helpers, but we keep the type
// permissive so the unmarshal won't fail on unexpected keys.
type indicatorDeleteRequest struct {
	RequestData platformTypes.DeleteIndicatorsRequest `json:"request_data"`
}

// --- InsertIndicators ------------------------------------------------------

func TestClient_InsertIndicators(t *testing.T) {
	t.Run("create returns added_objects with new rule_id", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+InsertIndicatorsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			var req indicatorInsertRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData, 1)
			assert.Equal(t, "evil.example.com", req.RequestData[0].Indicator)
			assert.Equal(t, enums.IndicatorTypeDomainName, req.RequestData[0].Type)
			assert.Equal(t, enums.IndicatorSeverityHigh, req.RequestData[0].Severity)
			// rule_id is omitempty and absent on create.
			assert.Zero(t, req.RequestData[0].RuleID)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [{"id": 42, "status": "Created a new indicator with the ID: 42 successfully"}],
				"updated_objects": [],
				"errors": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertIndicators(context.Background(), []platformTypes.Indicator{
			{
				Indicator: "evil.example.com",
				Type:      enums.IndicatorTypeDomainName,
				Severity:  enums.IndicatorSeverityHigh,
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

			var req indicatorInsertRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData, 1)
			assert.Equal(t, 42, req.RequestData[0].RuleID)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [],
				"updated_objects": [{"id": 42, "status": "Updated the indicator with the ID: 42 successfully"}],
				"errors": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertIndicators(context.Background(), []platformTypes.Indicator{
			{
				RuleID:    42,
				Indicator: "evil.example.com",
				Type:      enums.IndicatorTypeDomainName,
				Severity:  enums.IndicatorSeverityHigh,
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.UpdatedObjects, 1)
		assert.Equal(t, 42, resp.UpdatedObjects[0].ID)
		assert.Empty(t, resp.AddedObjects)
		assert.Empty(t, resp.Errors)
	})

	t.Run("errors are decoded as {index,status} objects (not strings)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"added_objects": [],
				"updated_objects": [],
				"errors": [
					{"index": 0, "status": "Failed to create indicator due to: Invalid IOC indicator"}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.InsertIndicators(context.Background(), []platformTypes.Indicator{{
			Indicator: "not-a-url",
			Type:      enums.IndicatorTypeURL,
			Severity:  enums.IndicatorSeverityLow,
		}})
		require.NoError(t, err)
		require.Len(t, resp.Errors, 1)
		assert.Equal(t, 0, resp.Errors[0].Index)
		assert.Contains(t, resp.Errors[0].Status, "Invalid IOC indicator")
		assert.Empty(t, resp.AddedObjects)
		assert.Empty(t, resp.UpdatedObjects)
	})

	t.Run("success body is read at top level (no reply wrapper)", func(t *testing.T) {
		// Regression guard for L3: the SDK should NOT expect a reply
		// wrapper on success. If a future refactor reintroduces
		// ResponseWrapperKeys: ["reply"] this test will fail because the
		// wrapped body below has no "reply" key.
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

		resp, err := client.InsertIndicators(context.Background(), []platformTypes.Indicator{{
			Indicator: "x", Type: enums.IndicatorTypeHash, Severity: enums.IndicatorSeverityInfo,
		}})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.AddedObjects[0].ID)
	})
}

// --- ListIndicators --------------------------------------------------------

func TestClient_ListIndicators(t *testing.T) {
	t.Run("returns objects with read-only fields populated", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListIndicatorsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req indicatorFilterRequest
			require.NoError(t, json.Unmarshal(body, &req))
			assert.True(t, req.RequestData.ExtendedView)
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "indicator", req.RequestData.Filters[0].Field)
			assert.Equal(t, "EQ", req.RequestData.Filters[0].Operator)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects_type": "indicator",
				"objects": [{
					"rule_id": 57,
					"indicator": "virus1.exe",
					"type": "FILENAME",
					"severity": "SEV_040_HIGH",
					"expiration_date": -1,
					"default_expiration_enabled": true,
					"comment": "test",
					"reputation": "BAD",
					"reliability": "C",
					"creation_time": 1781000000000,
					"modification_time": 1781000001000,
					"status": "ENABLED",
					"source": "Public API user (key #30)",
					"number_of_issues": 7
				}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListIndicators(context.Background(), platformTypes.ListIndicatorsRequest{
			ExtendedView: true,
			Filters: []platformTypes.IndicatorFilter{
				{Field: "indicator", Operator: "EQ", Value: "virus1.exe"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.ObjectsCount)
		assert.Equal(t, "indicator", resp.ObjectsType)
		require.Len(t, resp.Objects, 1)

		obj := resp.Objects[0]
		assert.Equal(t, 57, obj.RuleID)
		assert.Equal(t, "virus1.exe", obj.Indicator)
		assert.Equal(t, enums.IndicatorTypeFilename, obj.Type)
		assert.Equal(t, enums.IndicatorSeverityHigh, obj.Severity)
		assert.Equal(t, int64(-1), obj.ExpirationDate)
		assert.True(t, obj.DefaultExpirationEnabled)
		assert.Equal(t, "test", obj.Comment)
		assert.Equal(t, enums.IndicatorReputationBad, obj.Reputation)
		assert.Equal(t, enums.IndicatorReliabilityC, obj.Reliability)
		// Read-only fields (S3 fix):
		assert.Equal(t, int64(1781000000000), obj.CreationTime)
		assert.Equal(t, int64(1781000001000), obj.ModificationTime)
		assert.Equal(t, "ENABLED", obj.Status)
		assert.Equal(t, "Public API user (key #30)", obj.Source)
		assert.Equal(t, 7, obj.NumberOfIssues)
	})

	t.Run("handles JSON null reliability without error", func(t *testing.T) {
		// Live tenant returns "reliability": null when unset. Go's
		// encoding/json should map that to the zero value of the
		// string-aliased enum type ("").
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects_type": "indicator",
				"objects": [{
					"rule_id": 1,
					"indicator": "a",
					"type": "HASH",
					"severity": "SEV_010_INFO",
					"reliability": null
				}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListIndicators(context.Background(), platformTypes.ListIndicatorsRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Objects, 1)
		assert.Equal(t, enums.IndicatorReliability(""), resp.Objects[0].Reliability)
	})

	t.Run("empty result decodes cleanly", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": [], "objects_type": "indicator"}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListIndicators(context.Background(), platformTypes.ListIndicatorsRequest{})
		require.NoError(t, err)
		assert.Zero(t, resp.ObjectsCount)
		assert.Empty(t, resp.Objects)
	})
}

// --- DeleteIndicators ------------------------------------------------------

func TestClient_DeleteIndicators(t *testing.T) {
	t.Run("returns deleted rule_ids", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+DeleteIndicatorsEndpoint, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req indicatorDeleteRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "indicator", req.RequestData.Filters[0].Field)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 2, "objects": [101, 102]}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		ids, err := client.DeleteIndicators(context.Background(), platformTypes.DeleteIndicatorsRequest{
			Filters: []platformTypes.IndicatorFilter{
				{Field: "indicator", Operator: "EQ", Value: "evil.example.com"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []int{101, 102}, ids)
	})

	t.Run("idempotent: empty filter result is not an error", func(t *testing.T) {
		// Regression guard for C1: prior to the fix, DeleteIndicators
		// inspected a non-existent `success` field which was always
		// false, making this case look like a failure. Now the slice
		// length alone is the signal.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": []}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		ids, err := client.DeleteIndicators(context.Background(), platformTypes.DeleteIndicatorsRequest{
			Filters: []platformTypes.IndicatorFilter{
				{Field: "indicator", Operator: "EQ", Value: "does-not-exist"},
			},
		})
		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}

// --- FindIndicatorByName ---------------------------------------------------

func TestClient_FindIndicatorByName(t *testing.T) {
	t.Run("returns matching record", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var req indicatorFilterRequest
			require.NoError(t, json.Unmarshal(body, &req))
			require.Len(t, req.RequestData.Filters, 1)
			assert.Equal(t, "indicator", req.RequestData.Filters[0].Field)
			assert.Equal(t, "EQ", req.RequestData.Filters[0].Operator)
			assert.Equal(t, "evil.example.com", req.RequestData.Filters[0].Value)
			assert.True(t, req.RequestData.ExtendedView)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects": [{"rule_id": 9, "indicator": "evil.example.com", "type": "DOMAIN_NAME", "severity": "SEV_020_LOW"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindIndicatorByName(context.Background(), "evil.example.com")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 9, got.RuleID)
		assert.Equal(t, "evil.example.com", got.Indicator)
	})

	t.Run("returns nil when no match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": [], "objects_type": "indicator"}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindIndicatorByName(context.Background(), "missing")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("filters out value mismatches (defensive)", func(t *testing.T) {
		// The filter is server-evaluated, so in practice this won't fire,
		// but FindIndicatorByName double-checks the returned object's
		// `indicator` field. This test pins that contract.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects": [{"rule_id": 1, "indicator": "different-value", "type": "HASH", "severity": "SEV_020_LOW"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindIndicatorByName(context.Background(), "queried-value")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

// --- FindIndicatorByID -----------------------------------------------------

func TestClient_FindIndicatorByID(t *testing.T) {
	t.Run("queries by rule_id and returns match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			// The Value carries the int as JSON number; decode into a
			// loose struct to verify the on-wire type.
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
			// JSON numbers decode to float64 in an `any` field; the
			// integer payload of 192 round-trips losslessly.
			assert.Equal(t, float64(192), raw.RequestData.Filters[0].Value)
			assert.True(t, raw.RequestData.ExtendedView)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"objects_count": 1,
				"objects": [{"rule_id": 192, "indicator": "x", "type": "URL", "severity": "SEV_050_CRITICAL"}]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindIndicatorByID(context.Background(), 192)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 192, got.RuleID)
		assert.Equal(t, enums.IndicatorTypeURL, got.Type)
		assert.Equal(t, enums.IndicatorSeverityCritical, got.Severity)
	})

	t.Run("returns nil when no match", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"objects_count": 0, "objects": []}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		got, err := client.FindIndicatorByID(context.Background(), 999)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
