// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
)

// InsertIndicators uploads one or more IOCs. The endpoint upserts: submit
// a payload with no `rule_id` to create a new record, or include the
// server-assigned `rule_id` for an existing record to overwrite it.
// Per-record failures surface in Errors[] keyed by the original batch
// position; successful creates/updates appear in AddedObjects/UpdatedObjects.
//
// Note: the three indicator endpoints return success bodies at the top
// level (no `reply` wrapper). Error bodies do come wrapped in `reply` and
// are handled by the internal client's CortexCloudAPIError path before
// reaching response unmarshal.
func (c *Client) InsertIndicators(ctx context.Context, indicators []types.Indicator) (types.InsertIndicatorsResponse, error) {
	var resp types.InsertIndicatorsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, InsertIndicatorsEndpoint, nil, nil, indicators, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp, mapError(err)
}

// ListIndicators retrieves IOCs matching the provided filter set. Pass an
// empty Filters slice to list all indicators (capped by the server).
func (c *Client) ListIndicators(ctx context.Context, req types.ListIndicatorsRequest) (types.ListIndicatorsResponse, error) {
	var resp types.ListIndicatorsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListIndicatorsEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp, mapError(err)
}

// DeleteIndicators removes IOCs matching the provided filter set and
// returns the server-assigned rule_ids of the removed records. The API
// has no by-ID delete; identity is carried in the filter body (typically a
// single `{field:"indicator",operator:"EQ",value:<id>}` clause). An empty
// returned slice means the filter matched nothing — callers that treat
// "delete-by-name on a missing record" as success should ignore the count.
func (c *Client) DeleteIndicators(ctx context.Context, req types.DeleteIndicatorsRequest) ([]int, error) {
	var resp types.DeleteIndicatorsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteIndicatorsEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp.Objects, mapError(err)
}

// FindIndicatorByID returns the indicator whose server-assigned `rule_id`
// equals the given value, or nil if no match was found. The OpenAPI lists
// `rule_id` only in the indicators/get *example*, not in the documented
// filter-field enum, but the live API accepts it on EQ. Use it when you
// have the numeric ID but not the indicator string (e.g. cross-referencing
// audit logs).
func (c *Client) FindIndicatorByID(ctx context.Context, ruleID int) (*types.Indicator, error) {
	resp, err := c.ListIndicators(ctx, types.ListIndicatorsRequest{
		ExtendedView: true,
		Filters: []types.IndicatorFilter{
			{Field: "rule_id", Operator: "EQ", Value: ruleID},
		},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.Objects {
		if resp.Objects[i].RuleID == ruleID {
			return &resp.Objects[i], nil
		}
	}
	return nil, nil
}

// FindIndicatorByName is a convenience helper that returns the single
// indicator whose `indicator` field equals the given value, or nil if no
// match was found. It is the primary lookup path for Terraform-style
// resources where the indicator string is the identity.
func (c *Client) FindIndicatorByName(ctx context.Context, indicator string) (*types.Indicator, error) {
	resp, err := c.ListIndicators(ctx, types.ListIndicatorsRequest{
		ExtendedView: true,
		Filters: []types.IndicatorFilter{
			{Field: "indicator", Operator: "EQ", Value: indicator},
		},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.Objects {
		if resp.Objects[i].Indicator == indicator {
			return &resp.Objects[i], nil
		}
	}
	return nil, nil
}

