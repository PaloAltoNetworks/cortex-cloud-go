package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateAssetGroup(t *testing.T) {
	t.Run("should create asset group successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/%s", CreateAssetGroupEndpoint), r.URL.Path)

			var req map[string]types.CreateOrUpdateAssetGroupRequestWrapper
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "New Group", req["request_data"].AssetGroup.GroupName)
			assert.Equal(t, "Dynamic", req["request_data"].AssetGroup.GroupType)
			assert.Equal(t, "Description for New Group", req["request_data"].AssetGroup.GroupDescription)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{"data":{"success":true,"asset_group_id":123}}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		createReq := types.CreateOrUpdateAssetGroupRequest{
			GroupName: "New Group",
			GroupType: "Dynamic",
			GroupDescription: "Description for New Group",
			MembershipPredicate: types.CriteriaFilter{
				And: []types.Criteria{
					{
						SearchField: "xdm.asset.name",
						SearchType: "NCONTAINS",
						SearchValue: "test",
					},
				},
			},
		}

		resp, err := client.CreateAssetGroup(context.Background(), createReq)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, 123, resp.AssetGroupID)
	})
}

func TestClient_UpdateAssetGroup(t *testing.T) {
	t.Run("should update asset group successfully", func(t *testing.T) {
		requestURL := fmt.Sprintf("/%s123", UpdateAssetGroupEndpoint)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, requestURL, r.URL.Path)

			var req map[string]types.CreateOrUpdateAssetGroupRequestWrapper
			body, readBodyErr := io.ReadAll(r.Body)
			require.NoError(t, readBodyErr)

			unmarshalErr := json.Unmarshal(body, &req)
			require.NoError(t, unmarshalErr)

			assert.Equal(t, "Updated Name", req["request_data"].AssetGroup.GroupName)
			assert.Equal(t, "Updated Description", req["request_data"].AssetGroup.GroupDescription)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{"data":{"success":true}}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		updateReq := types.CreateOrUpdateAssetGroupRequest{
			GroupName: "Updated Name",
			GroupDescription: "Updated Description",
		}

		resp, err := client.UpdateAssetGroup(context.Background(), 123, updateReq)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})
}

func TestClient_DeleteAssetGroup(t *testing.T) {
	t.Run("should delete asset group successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/%s123", DeleteAssetGroupEndpoint), r.URL.Path) 

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{"data":{"success":true}}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.DeleteAssetGroup(context.Background(), 123)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})
}

func TestClient_ListAssetGroups(t *testing.T) {
	t.Run("should list asset groups successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListAssetGroupsEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": {
					"data": [
						{
							"XDM.ASSET_GROUP.ID": 1,
							"XDM.ASSET_GROUP.NAME": "grp_pcs_lab"
						}
					]
				}
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListAssetGroups(context.Background(), types.ListAssetGroupsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, 1, resp.Data[0].ID)
		assert.Equal(t, "grp_pcs_lab", resp.Data[0].Name)
	})
}
