// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
)

// GetUser retrieves the specified user in your environment.
func (c *Client) GetUser(ctx context.Context, userEmail string) (types.User, error) {
	var ans types.User
	resp, err := c.ListUsers(ctx)
	if err != nil {
		return ans, err
	}
	for _, user := range resp {
		if user.Email == userEmail {
			ans = user
			return ans, nil
		}
	}
	return ans, fmt.Errorf("no user found with email \"%s\"", userEmail)
}

// ListUsers retrieves a list of the current users in your environment.
func (c *Client) ListUsers(ctx context.Context) ([]types.User, error) {
	var ans []types.User
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListUsersEndpoint, nil, nil, nil, &ans, &client.DoOptions{
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

func (c *Client) ListRoles(ctx context.Context, roleNames []string) ([]types.Role, error) {
	var (
		resp            [][]types.Role
		normalizedRoles []types.Role
	)
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRolesEndpoint, nil, nil, roleNames, &resp, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data", "role_names"},
		ResponseWrapperKeys: []string{"reply"},
	})
	if err != nil {
		return []types.Role{}, err
	}
	if len(resp) == 0 {
		return []types.Role{}, fmt.Errorf("no roles found for provided value(s): \"%s\"", strings.Join(roleNames, "\", \""))
	}
	for _, innerSlice := range resp {
		if len(innerSlice) == 1 {
			normalizedRoles = append(normalizedRoles, innerSlice[0])
		}
	}
	return normalizedRoles, err
}

// SetRole adds or removes one or more users from a role.
//
// If no RoleName is provided in the SetRoleRequest, the user is removed from a role.
func (c *Client) SetRole(ctx context.Context, input types.SetRoleRequest) (types.SetRoleResponse, error) {
	var ans types.SetRoleResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, SetUserRoleEndpoint, nil, nil, input, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// GetRiskScore retrieves the risk score of a specific user or endpoint in your environment,
// along with the reason for the score.
func (c *Client) GetRiskScore(ctx context.Context, req types.GetRiskScoreRequest) (types.GetRiskScoreResponse, error) {
	var ans types.GetRiskScoreResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetRiskScoreEndpoint, nil, nil, req, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})

	return ans, err
}

// ListRiskyUsers retrieves a list of users with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyUsers(ctx context.Context) ([]types.ListRiskyUsersResponse, error) {
	var ans []types.ListRiskyUsersResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyUsersEndpoint, nil, nil, nil, &ans, &client.DoOptions{
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// ListRiskyHosts retrieves a list of endpoints with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyHosts(ctx context.Context) ([]types.ListRiskyHostsResponse, error) {
	var ans []types.ListRiskyHostsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyHostsEndpoint, nil, nil, nil, &ans, &client.DoOptions{
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// HealthCheck performs a health check on the service.
func (c *Client) HealthCheck(ctx context.Context) (types.HealthCheckResponse, error) {
	var ans types.HealthCheckResponse
	_, err := c.internalClient.Do(ctx, http.MethodGet, HealthCheckEndpoint, nil, nil, nil, &ans, &client.DoOptions{
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// GetTenantInfo retrieves information about the specified tenants.
func (c *Client) GetTenantInfo(ctx context.Context, req types.GetTenantInfoRequest) ([]types.TenantInfo, error) {
	var ans []types.TenantInfo
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetTenantInfoEndpoint, nil, nil, req, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply"},
	})
	return ans, err
}

// ListUserGroups retrieves a list of all user groups.
func (c *Client) ListUserGroups(ctx context.Context) ([]types.UserGroup, error) {
	var ans []types.UserGroup
	_, err := c.internalClient.Do(ctx, http.MethodGet, UserGroupEndpoint, nil, nil, nil, &ans, &client.DoOptions{
		ResponseWrapperKeys: []string{"data"},
	})
	return ans, err
}

// GetUserGroup retrieves information about the specified user groups.
func (c *Client) GetUserGroup(ctx context.Context, req types.GetUserGroupRequest) ([]types.UserGroup, error) {
	var ans []types.UserGroup
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetUserGroupEndpoint, nil, nil, req, &ans, &client.DoOptions{
		RequestWrapperKeys:  []string{"request_data"},
		ResponseWrapperKeys: []string{"reply", "data"},
	})
	return ans, err
}

// CreateUserGroup creates a new user group and returns its ID.
func (c *Client) CreateUserGroup(ctx context.Context, req types.UserGroupCreateRequest) (string, error) {
	var resp struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	_, err := c.internalClient.Do(ctx, http.MethodPost, UserGroupEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	if err != nil {
		return "", err
	}

	// Message is "user group with group id <id> created successfully"
	parts := strings.Split(resp.Data.Message, " ")
	if len(parts) < 6 || parts[0] != "user" || parts[1] != "group" {
		return "", fmt.Errorf("unexpected create user group response message: %s", resp.Data.Message)
	}
	groupID := parts[5]
	return groupID, nil
}

// EditUserGroup edits an existing user group.
// It takes a groupID and a UserGroupEditRequest object containing the fields to update.
func (c *Client) EditUserGroup(ctx context.Context, groupID string, req types.UserGroupEditRequest) (string, error) {
	var resp struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	// The request body is wrapped in {"request_data": ...} as seen in other API calls.
	_, err := c.internalClient.Do(ctx, http.MethodPatch, UserGroupEndpoint, &[]string{groupID}, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp.Data.Message, err
}

// DeleteUserGroup deletes an existing user group by its ID.
func (c *Client) DeleteUserGroup(ctx context.Context, groupID string) (string, error) {
	var resp struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	_, err := c.internalClient.Do(ctx, http.MethodDelete, UserGroupEndpoint, &[]string{groupID}, nil, nil, &resp, nil)
	return resp.Data.Message, err
}

// ListIAMUsers retrieves a list of all users and their respective properties.
func (c *Client) ListIAMUsers(ctx context.Context) (*types.ListIamUsersResponse, error) {
	var ans types.ListIamUsersResponse
	_, err := c.internalClient.Do(ctx, http.MethodGet, IamUsersEndpoint, nil, nil, nil, &ans, nil)
	if err != nil {
		return nil, err
	}
	return &ans, nil
}

// GetIAMUser retrieves a user and its respective properties.
func (c *Client) GetIAMUser(ctx context.Context, userEmail string) (*types.IamUser, error) {
	var ans types.GetIamUserResponse
	if _, err := c.internalClient.Do(ctx, http.MethodGet, IamUsersEndpoint, &[]string{userEmail}, nil, nil, &ans, nil); err != nil {
		return nil, err
	}
	return &ans.Data, nil
}

// EditIAMUser edits an existing user.
func (c *Client) EditIAMUser(ctx context.Context, userEmail string, req types.IamUserEditRequest) (map[string]any, error) {
	var resp map[string]any
	_, err := c.internalClient.Do(ctx, http.MethodPatch, IamUsersEndpoint, &[]string{userEmail}, nil, req, &resp, nil)
	return resp, err
}
