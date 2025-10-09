// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	cloudOnboardingTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/cloudonboarding"
	filterTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateOutpostTemplate(t *testing.T) {
	t.Run("should create outpost template successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/%s", CreateOutpostTemplateEndpoint), r.URL.Path)

			var reqBody struct {
				RequestData cloudOnboardingTypes.CreateOutpostTemplateRequest `json:"request_data"`
			}
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			assert.Equal(t, "AWS", reqBody.RequestData.CloudProvider)
			require.Len(t, reqBody.RequestData.CustomResourceTags, 1)
			assert.Equal(t, "key", reqBody.RequestData.CustomResourceTags[0].Key)
			assert.Equal(t, "value", reqBody.RequestData.CustomResourceTags[0].Value)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{"automated":{"link":"http://example.com","tracking_guid":"123"},"manual":{"CF":"manual-link"}}}`)
		})

		client, server := setupTest(t, handler)
		defer server.Close()

		input := cloudOnboardingTypes.CreateOutpostTemplateRequest{
			CloudProvider: "AWS",
			CustomResourceTags: []cloudOnboardingTypes.Tag{
				{Key: "key", Value: "value"},
			},
		}

		resp, err := client.CreateOutpostTemplate(context.Background(), input)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "http://example.com", resp.Automated.Link)
		assert.Equal(t, "123", resp.Automated.TrackingGuid)
		assert.Equal(t, "manual-link", resp.Manual.CF)
	})
}

func TestClient_UpdateOutpost(t *testing.T) {
	t.Run("should update outpost successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/%s", UpdateOutpostEndpoint), r.URL.Path)

			var reqBody struct {
				RequestData cloudOnboardingTypes.UpdateOutpostRequest `json:"request_data"`
			}
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			assert.Equal(t, "outpost-123", reqBody.RequestData.OutpostID)
			assert.Equal(t, "AWS", reqBody.RequestData.CloudProvider)
			require.Len(t, reqBody.RequestData.CustomResourceTags, 1)
			assert.Equal(t, "new-key", reqBody.RequestData.CustomResourceTags[0].Key)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{}}`)
		})

		client, server := setupTest(t, handler)
		defer server.Close()

		input := cloudOnboardingTypes.UpdateOutpostRequest{
			OutpostID:     "outpost-123",
			CloudProvider: "AWS",
			CustomResourceTags: []cloudOnboardingTypes.Tag{
				{Key: "new-key", Value: "new-value"},
			},
		}

		err := client.UpdateOutpost(context.Background(), input)
		require.NoError(t, err)
	})
}

func TestClient_ListOutposts(t *testing.T) {
	t.Run("should list outposts successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/%s", ListOutpostsEndpoint), r.URL.Path)

			var reqBody struct {
				RequestData cloudOnboardingTypes.ListOutpostsRequest `json:"request_data"`
			}
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			assert.NotNil(t, reqBody.RequestData.FilterData)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": {
					"DATA": [
						{
							"cloud_provider": "AWS",
							"outpost_id": "outpost-1",
							"created_at": 1740307555361,
							"type": "MANAGED"
						}
					],
					"FILTER_COUNT": 1,
					"TOTAL_COUNT": 1
				}
			}`)
		})

		client, server := setupTest(t, handler)
		defer server.Close()

		input := cloudOnboardingTypes.ListOutpostsRequest{
			FilterData: filterTypes.FilterData{},
		}

		resp, err := client.ListOutposts(context.Background(), input)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 1, resp.FilterCount)
		assert.Equal(t, 1, resp.TotalCount)
		require.Len(t, resp.Data, 1)
		assert.Equal(t, "AWS", resp.Data[0].CloudProvider)
		assert.Equal(t, "outpost-1", resp.Data[0].OutpostID)
	})
}
