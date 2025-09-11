// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/api"
	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	config := &api.Config{
		ApiUrl:    server.URL,
		ApiKey:    "test-key",
		ApiKeyId:  123,
		Transport: server.Client().Transport.(*http.Transport),
	}
	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	return client, server
}

func TestNewClient(t *testing.T) {
	t.Run("should return error for nil config", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Error(t, err)
		assert.NotNil(t, client)
		assert.Nil(t, client.internalClient)
		assert.Equal(t, "received nil api.Config", err.Error())
	})

	t.Run("should create new client with valid config", func(t *testing.T) {
		config := &api.Config{
			ApiUrl:   "https://api.example.com",
			ApiKey:   "test-key",
			ApiKeyId: 123,
		}
		client, err := NewClient(config)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.internalClient)
	})
}

func TestClient_ListPoliciesCompliance(t *testing.T) {
	t.Run("should return no policies", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/"+ListCloudWorkloadPoliciesEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"policies":[]}`) // Empty response
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		listReq := ListPoliciesRequest{
			PolicyTypes: []enums.PolicyType{enums.PolicyTypeCompliance},
		}
		resp, err := client.ListPolicies(context.Background(), listReq)
		assert.NoError(t, err)
		assert.Len(t, resp.Policies, 0)
	})
}

func TestClient_ListPoliciesMalware(t *testing.T) {
	t.Run("should list policies successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/"+ListCloudWorkloadPoliciesEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"policies":[{"id":"policy123","name":"Test Policy", "type":"MALWARE"}]}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		listReq := ListPoliciesRequest{
			PolicyTypes: []enums.PolicyType{enums.PolicyTypeMalware},
		}
		resp, err := client.ListPolicies(context.Background(), listReq)
		assert.NoError(t, err)
		assert.Len(t, resp.Policies, 1)
		assert.Equal(t, "policy123", resp.Policies[0].Id)
		assert.Equal(t, enums.PolicyTypeMalware, resp.Policies[0].Type)
	})
}

func TestClient_CreatePolicy(t *testing.T) {
	t.Run("should create policy successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+CreateCloudWorkloadPolicyEndpoint, r.URL.Path)

			var req CreatePolicyRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "New Policy", req.Data.Name)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":"newPolicy123"}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		createReq := CreatePolicyRequest{
			Data: PolicyData{
				Name: "New Policy",
			},
		}
		resp, err := client.CreatePolicy(context.Background(), createReq)
		assert.NoError(t, err)
		assert.Equal(t, "newPolicy123", resp.Id)
	})
}

func TestClient_GetPolicyDetails(t *testing.T) {
	t.Run("should get policy details successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/"+GetCloudWorkloadPolicyDetailsEndpoint+"/policy123", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"policy":{"id":"policy123","name":"Test Policy"}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.GetPolicyDetails(context.Background(), "policy123")
		assert.NoError(t, err)
		assert.Equal(t, "policy123", resp.Policy.Id)
		assert.Equal(t, "Test Policy", resp.Policy.Name)
	})
}

func TestClient_UpdatePolicy(t *testing.T) {
	t.Run("should update policy successfully", func(t *testing.T) {
		callCount := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount == 1 {
				// First call - GetPolicyDetails
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"policy":{"id":"policy123","name":"Original Policy"}}`)
			} else {
				// Second call - UpdatePolicy
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(http.StatusOK)
			}
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		updateReq := UpdatePolicyRequest{
			Id: "policy123",
			Data: PolicyData{
				Name: "Updated Policy",
			},
		}
		err := client.UpdatePolicy(context.Background(), updateReq)
		assert.NoError(t, err)
	})
}

func TestClient_DeletePolicy(t *testing.T) {
	t.Run("should delete policy successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/"+DeleteCloudWorkloadPolicyEndpoint+"/policy123", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		deleteReq := DeletePolicyRequest{
			Id: "policy123",
		}
		err := client.DeletePolicy(context.Background(), deleteReq)
		assert.NoError(t, err)
	})
}
