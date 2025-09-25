// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
)

// ----------------------------------------------------------------------------
// Get Accounts
// ----------------------------------------------------------------------------

type GetCloudAccountsRequestData struct {
	InstanceId string     `json:"instance_id"`
	FilterData types.FilterData `json:"filter_data"`
}

type ListAccountsByInstanceResponseReply struct {
	Data        ListAccountsByInstanceResponseData `json:"DATA"`
	FilterCount int                                `json:"FILTER_COUNT"`
	TotalCount  int                                `json:"TOTAL_COUNT"`
}

type ListAccountsByInstanceResponseData struct {
	Status      string `json:"status"`
	AccountName string `json:"account_name"`
	AccountId   string `json:"account_id"`
	Environment string `json:"environment"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
}

func (c *Client) ListAccountsByInstance(ctx context.Context, input GetCloudAccountsRequestData) (ListAccountsByInstanceResponseReply, error) {
	var ans ListAccountsByInstanceResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListAccountsByInstanceEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// Enable or disable cloud accounts

type EnableDisableAccountsInInstancesRequestData struct {
	Ids        []string `json:"ids"`
	InstanceId string   `json:"instance_id"`
	Enable     bool     `json:"enable"`
}

type EnableDisableAccountsInInstancesResponseReply struct{}

func (c *Client) EnableAccountsInInstance(ctx context.Context, instanceId string, accountIds []string) (EnableDisableAccountsInInstancesResponseReply, error) {
	return c.enableDisableAccountsInInstance(ctx, instanceId, accountIds, true)
}

func (c *Client) DisableAccountsInInstance(ctx context.Context, instanceId string, accountIds []string) (EnableDisableAccountsInInstancesResponseReply, error) {
	return c.enableDisableAccountsInInstance(ctx, instanceId, accountIds, false)
}

func (c *Client) enableDisableAccountsInInstance(ctx context.Context, instanceId string, accountIds []string, enable bool) (EnableDisableAccountsInInstancesResponseReply, error) {
	req := EnableDisableAccountsInInstancesRequestData{
		InstanceId: instanceId,
		Ids:        accountIds,
		Enable:     enable,
	}

	var ans EnableDisableAccountsInInstancesResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableDisableAccountsInInstancesEndpoint, nil, nil, req, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}
