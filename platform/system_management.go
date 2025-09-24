// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform/types"
)

// GetUser retrieves the specified user in your environment.
func (c *Client) GetUser(ctx context.Context, input types.GetUserRequest) (types.User, error) {
	var ans types.User
	resp, err := c.ListUsers(ctx)
	if err != nil {
		return ans, err
	}
	for _, user := range resp {
		if user.UserEmail == input.Email {
			ans = user
			return ans, nil
		}
	}
	return ans, fmt.Errorf("User with email %s not found", input.Email)
}

// ListUsers retrieves a list of the current users in your environment.
func (c *Client) ListUsers(ctx context.Context) ([]types.User, error) {
	var ans []types.User
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListUsersEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

func (c *Client) ListRoles(ctx context.Context, input types.ListRolesRequestData) ([][]types.ListRolesResponseReply, error) {
	var ans [][]types.ListRolesResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRolesEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

// SetRole adds or removes one or more users from a role.
//
// If no RoleName is provided in the SetRoleRequest, the user is removed from a role.
func (c *Client) SetRole(ctx context.Context, input types.SetRoleRequestData) (types.SetRoleResponseReply, error) {
	var ans types.SetRoleResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, SetUserRoleEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

// GetRiskScore retrieves the risk score of a specific user or endpoint in your environment,
// along with the reason for the score.
func (c *Client) GetRiskScore(ctx context.Context, input types.GetRiskScoreRequestData) (types.GetRiskScoreResponseReply, error) {
	var ans types.GetRiskScoreResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetRiskScoreEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ListRiskyUsers retrieves a list of users with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyUsers(ctx context.Context) ([]types.ListRiskyUsersResponseReply, error) {
	var ans []types.ListRiskyUsersResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyUsersEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})
	return ans, err
}

// ListRiskyHosts retrieves a list of endpoints with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyHosts(ctx context.Context) ([]types.ListRiskyHostsResponseReply, error) {
	var ans []types.ListRiskyHostsResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyHostsEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})
	return ans, err
}
