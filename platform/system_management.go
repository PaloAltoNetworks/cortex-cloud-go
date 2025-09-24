// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
)

type User struct {
	UserEmail     string   `json:"user_email"`
	UserFirstName string   `json:"user_first_name"`
	UserLastName  string   `json:"user_last_name"`
	RoleName      string   `json:"role_name"`
	LastLoggedIn  int      `json:"last_logged_in"`
	UserType      string   `json:"user_type"`
	Groups        []string `json:"groups"`
	Scope         Scope    `json:"scope"`
}

type Scope struct {
	Endpoints   Endpoints   `json:"endpoints"`
	CasesIssues CasesIssues `json:"cases_issues"`
}

type Endpoints struct {
	EndpointGroups EndpointGroups `json:"endpoint_groups"`
	EndpointTags   EndpointTags   `json:"endpoint_tags"`
	Mode           string         `json:"mode"`
}

type EndpointGroups struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type EndpointTags struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type CasesIssues struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type Reason struct {
	DateCreated string `json:"date_created"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	Points      int    `json:"points"`
}

// ---------------------------
// Request/Response structs
// ---------------------------

// Get User

type GetUserRequest struct {
	Email string `json:"email"`
}

// List Roles

type ListRolesRequestData struct {
	// TODO: add validation tag/function for role names?
	RoleNames []string `json:"role_names" validate:"required,min=1"`
}

type ListRolesResponseReply struct {
	PrettyName  string   `json:"pretty_name"`
	Permissions []string `json:"permissions"`
	InsertTime  int      `json:"insert_time"`
	UpdateTime  int      `json:"update_time"`
	CreatedBy   string   `json:"created_by"`
	Description string   `json:"description"`
	Tags        string   `json:"tags"`
	Groups      []string `json:"groups"`
	Users       []string `json:"users"`
}

// Set Role

type SetRoleRequestData struct {
	UserEmails []string `json:"user_emails" validate:"required,min=1,dive,required,email"`
	RoleName   string   `json:"role_name"`
}

type SetRoleResponseReply struct {
	UpdateCount int `json:"update_count"`
}

// Get Risk Score

type GetRiskScoreRequestData struct {
	ID string `json:"id" validate:"required,sysmgmtID"`
}

type GetRiskScoreResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// List Risky Users

type ListRiskyUsersResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// List Risky Hosts

type ListRiskyHostsResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
}

// ---------------------------
// Request functions
// ---------------------------

// GetUser retrieves the specified user in your environment.
func (c *Client) GetUser(ctx context.Context, input GetUserRequest) (User, error) {

	var ans User

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
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {

	var ans []User
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListUsersEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

func (c *Client) ListRoles(ctx context.Context, input ListRolesRequestData) ([][]ListRolesResponseReply, error) {
	var ans [][]ListRolesResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRolesEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// SetRole adds or removes one or more users from a role.
//
// If no RoleName is provided in the SetRoleRequest, the user is removed from a role.
func (c *Client) SetRole(ctx context.Context, input SetRoleRequestData) (SetRoleResponseReply, error) {
	var ans SetRoleResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, SetUserRoleEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// GetRiskScore retrieves the risk score of a specific user or endpoint in your environment,
// along with the reason for the score.
func (c *Client) GetRiskScore(ctx context.Context, input GetRiskScoreRequestData) (GetRiskScoreResponseReply, error) {
	var ans GetRiskScoreResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetRiskScoreEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ListRiskyUsers retrieves a list of users with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyUsers(ctx context.Context) ([]ListRiskyUsersResponseReply, error) {
	var ans []ListRiskyUsersResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyUsersEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ListRiskyHosts retrieves a list of endpoints with the highest risk score in your environment
// along with the reason affecting each score.
func (c *Client) ListRiskyHosts(ctx context.Context) ([]ListRiskyHostsResponseReply, error) {
	var ans []ListRiskyHostsResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListRiskyHostsEndpoint, nil, nil, nil, &ans, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})

	return ans, err
}
