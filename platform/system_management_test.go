// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListUsers(t *testing.T) {
	t.Run("should list users successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListUsersEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": [
					{
						"user_email": "test@example.com",
						"user_first_name": "Test",
						"user_last_name": "User",
						"role_name": "Admin"
					}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListUsers(context.Background())
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "test@example.com", resp[0].UserEmail)
	})
}

func TestClient_ListRoles(t *testing.T) {
	t.Run("should list roles successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListRolesEndpoint, r.URL.Path)

			var req map[string]ListRolesRequestData
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, []string{"Admin", "User"}, req["request_data"].RoleNames)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": [
					[
					{
						"pretty_name": "Admin",
						"permissions": [
						"Reports",
						"Playbooks",
						"Datasets Access Control",
						"Dashboards",
						"Scripts"
						],
						"insert_time": 1658315576844,
						"update_time": null,
						"created_by": "admin@example.com",
						"description": "",
						"groups": [
						"group1",
						"group2"
						],
						"users": ["admin@example.com"]
					}
					],
					[
					{
						"pretty_name": "User",
						"permissions": [
						"Dashboards",
						"Datasets Access Control"
						],
						"insert_time": 1661435660656,
						"update_time": null,
						"created_by": "admin@example.com",
						"description": "",
						"groups": [],
						"users": ["user@example.com"]
					}
					]
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		listReq := ListRolesRequestData{
			RoleNames: []string{"Admin", "User"},
		}
		resp, err := client.ListRoles(context.Background(), listReq)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Len(t, resp[0], 1)
		assert.Len(t, resp[1], 1)
		assert.Equal(t, "Admin", resp[0][0].PrettyName)
		assert.Equal(t, "admin@example.com", resp[0][0].CreatedBy)
		assert.Len(t, resp[0][0].Users, 1)
		assert.Equal(t, "admin@example.com", resp[0][0].Users[0])
		assert.Equal(t, "User", resp[1][0].PrettyName)
		assert.Equal(t, "admin@example.com", resp[1][0].CreatedBy)
		assert.Len(t, resp[1][0].Users, 1)
		assert.Equal(t, "user@example.com", resp[1][0].Users[0])
	})
}

func TestClient_SetRole(t *testing.T) {
	t.Run("should set role successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+SetUserRoleEndpoint, r.URL.Path)

			var req map[string]SetRoleRequestData
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "new-role", req["request_data"].RoleName)
			assert.Equal(t, []string{"user@example.com"}, req["request_data"].UserEmails)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply": {"update_count": "1"}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		setReq := SetRoleRequestData{
			RoleName:   "new-role",
			UserEmails: []string{"user@example.com"},
		}
		resp, err := client.SetRole(context.Background(), setReq)
		assert.NoError(t, err)
		assert.Equal(t, "1", resp.UpdateCount)
	})
}

func TestClient_GetRiskScore(t *testing.T) {
	t.Run("should get risk score successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+GetRiskScoreEndpoint, r.URL.Path)

			var req map[string]GetRiskScoreRequestData
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "user123", req["request_data"].ID)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": {
					"id": "user123",
					"score": 95,
					"risk_level": "Critical"
				}
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		getReq := GetRiskScoreRequestData{
			ID: "user123",
		}
		resp, err := client.GetRiskScore(context.Background(), getReq)
		assert.NoError(t, err)
		assert.Equal(t, "user123", resp.ID)
		assert.Equal(t, 95, resp.Score)
	})
}

func TestClient_ListRiskyUsers(t *testing.T) {
	t.Run("should list risky users successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListRiskyUsersEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": [
					{
						"id": "user456",
						"email": "risky@example.com",
						"score": 80,
						"risk_level": "High"
					}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListRiskyUsers(context.Background())
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "risky@example.com", resp[0].Email)
	})
}

func TestClient_ListRiskyHosts(t *testing.T) {
	t.Run("should list risky hosts successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListRiskyHostsEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": [
					{
						"id": "host789",
						"score": 70,
						"risk_level": "Medium"
					}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListRiskyHosts(context.Background())
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "host789", resp[0].ID)
	})
}
