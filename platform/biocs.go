// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package platform

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	types "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
)

// InsertBIOCs uploads one or more BIOCs. The endpoint upserts: submit a
// payload with no `rule_id` to create a new record, or include the
// server-assigned `rule_id` for an existing record to overwrite it.
// Per-record failures surface in Errors[] keyed by the original batch
// position; successful creates/updates appear in AddedObjects/UpdatedObjects.
//
// Unlike /indicators/insert, /bioc/insert returns HTTP 400 (not 200) when
// any single record fails validation, but the body still uses the success
// shape (added_objects/updated_objects/errors). This helper recovers the
// typed response from that body and returns (resp, nil) so callers can
// inspect resp.Errors without first handling an HTTP error. True transport-
// level failures (no body, malformed body, auth, etc.) still surface as a
// non-nil error from mapError. Verified against a live tenant.
func (c *Client) InsertBIOCs(ctx context.Context, biocs []types.BIOC) (types.InsertBIOCsResponse, error) {
	var resp types.InsertBIOCsResponse
	body, err := c.internalClient.Do(ctx, http.MethodPost, InsertBIOCsEndpoint, nil, nil, biocs, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	if err == nil {
		return resp, nil
	}
	// Recover the success-shape body that /bioc/insert returns on per-
	// record validation failures. If the body parses and contains any of
	// the expected fields, surface it instead of the HTTP-level error.
	var recovered types.InsertBIOCsResponse
	if jsonErr := json.Unmarshal(body, &recovered); jsonErr == nil {
		if len(recovered.Errors) > 0 || len(recovered.AddedObjects) > 0 || len(recovered.UpdatedObjects) > 0 {
			return recovered, nil
		}
	}
	return resp, mapError(err)
}

// ListBIOCs retrieves BIOCs matching the provided filter set. Pass an empty
// Filters slice to list all BIOCs (capped by the server).
func (c *Client) ListBIOCs(ctx context.Context, req types.ListBIOCsRequest) (types.ListBIOCsResponse, error) {
	var resp types.ListBIOCsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListBIOCsEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp, mapError(err)
}

// DeleteBIOCs removes BIOCs matching the provided filter set and returns
// the server-assigned rule_ids of the removed records. The API has no by-ID
// delete; identity is carried in the filter body. Callers managing a single
// BIOC by rule_id should pass `{field:"rule_id",operator:"EQ",value:<id>}`
// — `rule_id` is undocumented in the OpenAPI filter enum but accepted by
// the live API (verified against a live tenant). Deleting by
// `name` is unsafe because BIOC names are not unique per tenant.
//
// An empty returned slice means the filter matched nothing — callers that
// treat "delete on a missing record" as success should ignore the count.
func (c *Client) DeleteBIOCs(ctx context.Context, req types.DeleteBIOCsRequest) ([]int, error) {
	var resp types.DeleteBIOCsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteBIOCsEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKeys: []string{"request_data"},
	})
	return resp.Objects, mapError(err)
}

// FindBIOCByID returns the BIOC whose server-assigned `rule_id` equals the
// given value, or nil if no match was found. The OpenAPI lists `rule_id`
// only in the bioc/get *example*, not in the documented filter-field enum,
// but the live API accepts it on EQ. Use this as the primary lookup path —
// BIOC names are not unique per tenant, so FindBIOCByName is unsafe for
// stateful flows.
func (c *Client) FindBIOCByID(ctx context.Context, ruleID int) (*types.BIOC, error) {
	resp, err := c.ListBIOCs(ctx, types.ListBIOCsRequest{
		ExtendedView: true,
		Filters: []types.BIOCFilter{
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

// FindBIOCByName returns the FIRST BIOC whose `name` equals the given
// value, or nil if no match was found. BIOC names are not unique on a
// tenant: two BIOCs may share the same name with distinct rule_ids. This
// helper is useful for ad-hoc CLI lookups but unsafe for stateful flows —
// use FindBIOCByID once you know the rule_id.
func (c *Client) FindBIOCByName(ctx context.Context, name string) (*types.BIOC, error) {
	resp, err := c.ListBIOCs(ctx, types.ListBIOCsRequest{
		ExtendedView: true,
		Filters: []types.BIOCFilter{
			{Field: "name", Operator: "EQ", Value: name},
		},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.Objects {
		if resp.Objects[i].Name == name {
			return &resp.Objects[i], nil
		}
	}
	return nil, nil
}
