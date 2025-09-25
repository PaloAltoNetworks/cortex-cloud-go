// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package appsec

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"runtime"
//	"testing"
//
//	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestBuildInfo(t *testing.T) {
//	if os.Getenv("CI") == "" {
//		t.Skip("Skipping build info test on local machine.")
//	}
//	expectedGitCommit := "test123"
//	expectedGoVersion := runtime.Version()
//	expectedBuildDate := "0000-00-00T00:00:00+0000"
//
//	t.Run("should return expected build info", func(t *testing.T) {
//		assert.Equal(t, expectedGitCommit, GitCommit)
//		assert.Equal(t, expectedGoVersion, GoVersion)
//		assert.Equal(t, expectedBuildDate, BuildDate)
//	})
//}
//
//func TestNewClient(t *testing.T) {
//	t.Run("should return error for nil config", func(t *testing.T) {
//		client, err := NewClient(nil)
//		assert.Error(t, err)
//		assert.NotNil(t, client)
//		assert.Nil(t, client.internalClient)
//		assert.Equal(t, "received nil api.Config", err.Error())
//	})
//
//	t.Run("should create new client with valid config", func(t *testing.T) {
//		config := &api.Config{
//			ApiUrl:   "https://api.example.com",
//			ApiKey:   "test-key",
//			ApiKeyId: 123,
//		}
//		client, err := NewClient(config)
//		assert.NoError(t, err)
//		assert.NotNil(t, client)
//		assert.NotNil(t, client.internalClient)
//	})
//}
//
//func setupTest(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
//	server := httptest.NewServer(handler)
//	config := &api.Config{
//		ApiUrl:    server.URL,
//		ApiKey:    "test-key",
//		ApiKeyId:  123,
//		Transport: server.Client().Transport.(*http.Transport),
//	}
//	client, err := NewClient(config)
//	assert.NoError(t, err)
//	assert.NotNil(t, client)
//	return client, server
//}
//
//func TestClient_Get(t *testing.T) {
//	t.Run("should get a rule successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodGet, r.Method)
//			assert.Equal(t, "/public_api/appsec/v1/rules/rule123", r.URL.Path)
//			w.WriteHeader(http.StatusOK)
//			fmt.Fprint(w, `{"id":"rule123","name":"Test Rule"}`)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		rule, err := client.Get(context.Background(), "rule123")
//		assert.NoError(t, err)
//		assert.Equal(t, "rule123", rule.Id)
//		assert.Equal(t, "Test Rule", rule.Name)
//	})
//}
//
//func TestClient_List(t *testing.T) {
//	t.Run("should list rules successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodGet, r.Method)
//			assert.Equal(t, "/"+RulesEndpoint, r.URL.Path)
//			assert.Equal(t, "true", r.URL.Query().Get("enabled"))
//			assert.Equal(t, "10", r.URL.Query().Get("limit"))
//			w.WriteHeader(http.StatusOK)
//			fmt.Fprint(w, `{"rules":[{"id":"rule123","name":"Test Rule"}],"offset":0}`)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		listReq := ListRequest{
//			Enabled: true,
//			Limit:   10,
//		}
//		resp, err := client.List(context.Background(), listReq)
//		assert.NoError(t, err)
//		assert.Len(t, resp.Rules, 1)
//		assert.Equal(t, "rule123", resp.Rules[0].Id)
//	})
//}
//
//func TestClient_CreateOrClone(t *testing.T) {
//	t.Run("should create rule successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodPost, r.Method)
//			assert.Equal(t, "/"+RulesEndpoint, r.URL.Path)
//
//			var req CreateOrCloneRequest
//			err := json.NewDecoder(r.Body).Decode(&req)
//			assert.NoError(t, err)
//			assert.Equal(t, "New Rule", req.Name)
//
//			w.WriteHeader(http.StatusCreated)
//			fmt.Fprint(w, `{"id":"newRule123","name":"New Rule"}`)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		createReq := CreateOrCloneRequest{
//			Name: "New Rule",
//		}
//		rule, err := client.CreateOrClone(context.Background(), createReq)
//		assert.NoError(t, err)
//		assert.Equal(t, "newRule123", rule.Id)
//	})
//}
//
//// TODO: uncomment and fix after fixing Update endpoint behaviour
////func TestClient_Update(t *testing.T) {
////	t.Run("should update a rule successfully", func(t *testing.T) {
////		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
////			//assert.Equal(t, http.MethodPatch, r.Method)
////			assert.Equal(t, http.MethodGet, r.Method)
////			assert.Equal(t, fmt.Sprintf("/%s/rule123", RulesEndpoint), r.URL.Path)
////
////			var req UpdateRequest
////			err := json.NewDecoder(r.Body).Decode(&req)
////			assert.NoError(t, err)
////			assert.Equal(t, "Updated Name", req.Name)
////
////			w.WriteHeader(http.StatusOK)
////			fmt.Fprint(w, `{"rule":{"id":"rule123","name":"Updated Name"}}`)
////		})
////		client, server := setupTest(t, handler)
////		defer server.Close()
////
////		updateReq := UpdateRequest{
////			Name: "Updated Name",
////		}
////		resp, err := client.Update(context.Background(), "rule123", updateReq)
////		assert.NoError(t, err)
////		assert.Equal(t, "Updated Name", resp.Rule.Name)
////	})
////}
//
//func TestClient_Delete(t *testing.T) {
//	t.Run("should delete a rule successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodDelete, r.Method)
//			assert.Equal(t, fmt.Sprintf("/%s/rule123", RulesEndpoint), r.URL.Path)
//			w.WriteHeader(http.StatusNoContent)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		err := client.Delete(context.Background(), "rule123")
//		assert.NoError(t, err)
//	})
//}
//
//func TestClient_Validate(t *testing.T) {
//	t.Run("should validate a rule successfully", func(t *testing.T) {
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			assert.Equal(t, http.MethodPost, r.Method)
//			assert.Equal(t, "/"+RulesValidationEndpoint, r.URL.Path)
//
//			bodyBytes, err := io.ReadAll(r.Body)
//			assert.NoError(t, err)
//
//			var req []ValidateRequest
//			err = json.Unmarshal(bodyBytes, &req)
//			assert.NoError(t, err)
//			assert.Len(t, req, 1)
//			assert.Equal(t, "my-framework", req[0].Framework)
//
//			w.WriteHeader(http.StatusOK)
//			isValid := true
//			err = json.NewEncoder(w).Encode(ValidateResponse{IsValid: &isValid})
//			assert.NoError(t, err)
//		})
//		client, server := setupTest(t, handler)
//		defer server.Close()
//
//		validateReq := []ValidateRequest{
//			{
//				Framework:  "my-framework",
//				Definition: "my-definition",
//			},
//		}
//		resp, err := client.Validate(context.Background(), validateReq)
//		assert.NoError(t, err)
//		assert.NotNil(t, resp.IsValid)
//		assert.True(t, *resp.IsValid)
//	})
//}
