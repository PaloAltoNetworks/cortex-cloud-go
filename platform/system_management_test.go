// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		users, err := client.ListUsers(context.Background())
		assert.NoError(t, err)
		require.Len(t, users, 1)
		user := users[0]
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.Equal(t, "Admin", user.RoleName)
	})
}

func TestClient_ListRoles(t *testing.T) {
	t.Run("should list roles successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+ListRolesEndpoint, r.URL.Path)

			var req map[string]map[string][]string
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			require.NotNil(t, req["request_data"]["role_names"])
			assert.Equal(t, []string{"Admin", "User"}, req["request_data"]["role_names"])

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

		roles, err := client.ListRoles(context.Background(), []string{"Admin", "User"})
		assert.NoError(t, err)
		require.Len(t, roles, 2)
		assert.Equal(t, "Admin", roles[0].PrettyName)
		assert.Equal(t, "admin@example.com", roles[0].CreatedBy)
		assert.Len(t, roles[0].Users, 1)
		assert.Equal(t, "admin@example.com", roles[0].Users[0])
		assert.Equal(t, "User", roles[1].PrettyName)
		assert.Equal(t, "admin@example.com", roles[1].CreatedBy)
		assert.Len(t, roles[1].Users, 1)
		assert.Equal(t, "user@example.com", roles[1].Users[0])
	})
}

func TestClient_SetRole(t *testing.T) {
	t.Run("should set role successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+SetUserRoleEndpoint, r.URL.Path)

			var req map[string]types.SetRoleRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "new-role", req["request_data"].RoleName)
			assert.Equal(t, []string{"user@example.com"}, req["request_data"].UserEmails)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"reply": {"update_count": "1"}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		setReq := types.SetRoleRequest{
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

			var req map[string]types.GetRiskScoreRequest
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

		getReq := types.GetRiskScoreRequest{
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

func TestClient_HealthCheck(t *testing.T) {
	t.Run("should return health check status successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/"+HealthCheckEndpoint, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": {
					"service": "Cortex API",
					"status": "OK",
					"reason": "",
					"timestamp": 1678886400000
				}
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.HealthCheck(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "Cortex API", resp.Service)
		assert.Equal(t, "OK", resp.Status)
	})
}

func TestClient_GetTenantInfo(t *testing.T) {
	t.Run("should get tenant info successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+GetTenantInfoEndpoint, r.URL.Path)

			var req map[string]types.GetTenantInfoRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, []string{"tenant1"}, req["request_data"].Tenants)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"reply": [
					{
						"tenant_id": "tenant1",
						"tenant_name": "Tenant One"
					}
				]
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		getReq := types.GetTenantInfoRequest{
			Tenants: []string{"tenant1"},
		}
		resp, err := client.GetTenantInfo(context.Background(), getReq)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "tenant1", resp[0].TenantID)
		assert.Equal(t, "Tenant One", resp[0].TenantName)
	})
}

func TestClient_ListUserGroups(t *testing.T) {
	t.Run("should list user groups successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, UserGroupEndpoint, r.URL.Path)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": [
					{
						"group_id": "test_group1",
						"group_name": "Group1",
						"description": "Group One",
						"role_name": "role_name01",
						"pretty_role_name": "Role Name 01",
						"created_by": "user1@test.com",
						"updated_by": "user1@test.com",
						"created_ts": 1661170832341,
						"updated_ts": 1661171650679,
						"users": ["user1@test.com"],
						"group_type": "custom",
						"nested_groups": [],
						"idp_groups": []
					}
				],
				"metadata": {"total_count": 1}
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.ListUserGroups(context.Background())
		assert.NoError(t, err)
		require.Len(t, resp, 1)
		assert.Equal(t, "test_group1", resp[0].GroupID)
		assert.Equal(t, "Group1", resp[0].GroupName)
		assert.Equal(t, "Group One", resp[0].Description)
	})
}

func TestClient_GetUserGroup(t *testing.T) {
	t.Run("should get user group successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/"+GetUserGroupEndpoint, r.URL.Path)

			var req map[string]types.GetUserGroupRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, []string{"group1"}, req["request_data"].GroupNames)

			w.WriteHeader(http.StatusOK)
			// This mock response is based on the API spec provided by the user.
			fmt.Fprint(w, `{
				"reply": {
					"data": [
						{
							"group_id": "test_group1",
							"group_name": "Group1",
							"description": "Group One",
							"role_name": "role_name01",
							"pretty_role_name": "Role Name 01",
							"created_by": "user1@test.com",
							"updated_by": "user1@test.com",
							"created_ts": 1661170832341,
							"updated_ts": 1661171650679,
							"users": ["user1@test.com"],
							"group_type": "custom",
							"nested_groups": [],
							"idp_groups": []
						}
					],
					"metadata": {"total_count": 1}
				}
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		getReq := types.GetUserGroupRequest{
			GroupNames: []string{"group1"},
		}
		resp, err := client.GetUserGroup(context.Background(), getReq)
		assert.NoError(t, err)
		require.Len(t, resp, 1)
		assert.Equal(t, "test_group1", resp[0].GroupID)
		assert.Equal(t, "Group1", resp[0].GroupName)
		assert.Equal(t, "Group One", resp[0].Description)
	})
}

func TestClient_CreateUserGroup(t *testing.T) {
	t.Run("should create user group successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, UserGroupEndpoint, r.URL.Path)

			var req map[string]types.UserGroupCreateRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "New Group", req["request_data"].GroupName)
			assert.Equal(t, "A new group for testing", req["request_data"].Description)

			w.WriteHeader(http.StatusCreated)
			// A create operation typically returns the newly created object.
			fmt.Fprint(w, `{
				"group_id": "new-group-id",
				"group_name": "New Group",
				"description": "A new group for testing",
				"role_name": "role_name02",
				"pretty_role_name": "Role Name 02",
				"created_by": "creator@test.com",
				"updated_by": "creator@test.com",
				"created_ts": 1670000000000,
				"updated_ts": 1670000000000,
				"users": [],
				"group_type": "custom",
				"nested_groups": [],
				"idp_groups": []
			}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		createReq := types.UserGroupCreateRequest{
			GroupName:   "New Group",
			Description: "A new group for testing",
		}
		resp, err := client.CreateUserGroup(context.Background(), createReq)
		assert.NoError(t, err)
		assert.Equal(t, "new-group-id", resp.GroupID)
		assert.Equal(t, "New Group", resp.GroupName)
	})
}

func TestClient_EditUserGroup(t *testing.T) {
	t.Run("should edit user group successfully", func(t *testing.T) {
		const groupID = "group-to-edit-id"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, fmt.Sprintf(UserGroupEndpoint+"/%s", groupID), r.URL.Path)

			var req map[string]types.UserGroupEditRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			require.NotNil(t, req["request_data"].GroupName)

			w.WriteHeader(http.StatusOK)
			// A successful PATCH often returns a simple success message.
			fmt.Fprint(w, `{"reply": {"success": true}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		newName := "Updated Name"
		editReq := types.UserGroupEditRequest{
			GroupName: newName,
		}
		resp, err := client.EditUserGroup(context.Background(), groupID, editReq)
		assert.NoError(t, err)
		assert.True(t, resp["success"].(bool))
	})
}

func TestClient_DeleteUserGroup(t *testing.T) {
	t.Run("should delete user group successfully", func(t *testing.T) {
		const groupID = "group-to-delete-id"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, fmt.Sprintf(UserGroupEndpoint+"/%s", groupID), r.URL.Path)

			w.WriteHeader(http.StatusOK)
			// A successful DELETE often returns a simple success message.
			fmt.Fprint(w, `{"reply": {"success": true}}`)
		})
		client, server := setupTest(t, handler)
		defer server.Close()

		resp, err := client.DeleteUserGroup(context.Background(), groupID)
		assert.NoError(t, err)
		assert.True(t, resp["success"].(bool))
	})
}
