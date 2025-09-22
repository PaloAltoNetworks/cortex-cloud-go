package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateAssetGroup(t *testing.T) {
	t.Run("should create asset group successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+CreateAssetGroupEndpoint, r.URL.Path)

			var req map[string]CreateAssetGroupRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "New Group", req["asset_group"].GroupName)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply":{"data":{"success":true,"asset_group_id":123}}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		createReq := CreateAssetGroupRequest{
			GroupName: "New Group",
			GroupType: "dynamic",
		}
		// TODO: fix
		_, err := client.CreateAssetGroup(context.Background(), createReq)
		assert.NoError(t, err)
		//assert.True(t, resp.Success)
		//assert.Equal(t, 123, resp.AssetGroupID)
	})
}

//func TestClient_UpdateAssetGroup(t *testing.T) {
//	t.Run("should update asset group successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodPost, r.Method)
//			assert.Equal(t, "/"+UpdateAssetGroupEndpoint+"/group123", r.URL.Path)
//
//			var req map[string]UpdateAssetGroupRequest
//			err := json.NewDecoder(r.Body).Decode(&req)
//			assert.NoError(t, err)
//			assert.Equal(t, "Updated Group", req["asset_group"].GroupName)
//
//			w.WriteHeader(http.StatusOK)
//			fmt.Fprint(w, `{"reply":{"data":{"success":true}}}`)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		updateReq := UpdateAssetGroupRequest{
//			GroupName: "Updated Group",
//		}
//		resp, err := client.UpdateAssetGroup(context.Background(), "group123", updateReq)
//		assert.NoError(t, err)
//		assert.True(t, resp.Success)
//	})
//}

//func TestClient_DeleteAssetGroup(t *testing.T) {
//	t.Run("should delete asset group successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodPost, r.Method)
//			assert.Equal(t, "/"+DeleteAssetGroupEndpoint+"/group123", r.URL.Path)
//
//			w.WriteHeader(http.StatusOK)
//			fmt.Fprint(w, `{"reply":{"data":{"success":true}}}`)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		resp, err := client.DeleteAssetGroup(context.Background(), "group123")
//		assert.NoError(t, err)
//		assert.True(t, resp.Success)
//	})
//}

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

		resp, err := client.ListAssetGroups(context.Background(), ListAssetGroupsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, 1, resp.Data[0].ID)
		assert.Equal(t, "grp_pcs_lab", resp.Data[0].Name)
	})
}
