// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
)

// ----------------------------------------------------------------------------
// Indicator (IOC)
// ----------------------------------------------------------------------------

// Indicator describes a single IOC as accepted by /indicators/insert and
// returned by /indicators/get.
//
// The first nine fields are read+write — the live /indicators/insert
// endpoint accepts exactly this set (the API explicitly rejects unknown
// keys with a "The allowed fields are: {...}. Got also the fields: [...]"
// error). The remaining five fields are read-only: they appear in GET
// responses but are server-assigned and not part of the insert schema.
type Indicator struct {
	// --- Read + Write fields ---
	RuleID                   int                        `json:"rule_id,omitempty"`
	Indicator                string                     `json:"indicator"`
	Type                     enums.IndicatorType        `json:"type"`
	Severity                 enums.IndicatorSeverity    `json:"severity"`
	ExpirationDate           int64                      `json:"expiration_date,omitempty"`
	DefaultExpirationEnabled bool                       `json:"default_expiration_enabled,omitempty"`
	Comment                  string                     `json:"comment,omitempty"`
	Reputation               enums.IndicatorReputation  `json:"reputation,omitempty"`
	Reliability              enums.IndicatorReliability `json:"reliability,omitempty"`

	// --- Read-only fields (server-assigned, returned by /indicators/get) ---
	CreationTime     int64  `json:"creation_time,omitempty"`
	ModificationTime int64  `json:"modification_time,omitempty"`
	Status           string `json:"status,omitempty"`
	Source           string `json:"source,omitempty"`
	NumberOfIssues   int    `json:"number_of_issues,omitempty"`
}

// IndicatorFilter is a single filter expression used by indicators/get and
// indicators/delete. It is not the generic FilterRoot type because the
// indicator endpoints accept a flat list of {field, operator, value} entries
// (no nested AND/OR), and `value` is polymorphic — boolean fields like
// `default_expiration_enabled` MUST be sent as JSON booleans (sending a
// string returns HTTP 500); numeric fields like `expiration_date` MUST be
// integers; everything else passes through as a string.
type IndicatorFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// The Cortex `/indicators/insert` endpoint upserts: submit with no
// `rule_id` to create, or include the existing `rule_id` to overwrite the
// matching record. The SDK takes a []Indicator slice directly and the
// internal client wraps it as `{"request_data": [...]}` — there is no
// dedicated request struct.

// ListIndicatorsRequest is the request body for indicators/get. ExtendedView
// toggles inclusion of vendor and class fields in the response — though in
// practice the live API ignores it and returns the same field set either
// way (vendor/class are not populated for IOCs created via /insert).
type ListIndicatorsRequest struct {
	Filters      []IndicatorFilter `json:"filters,omitempty"`
	ExtendedView bool              `json:"extended_view,omitempty"`
	SearchFrom   int               `json:"search_from,omitempty"`
	SearchTo     int               `json:"search_to,omitempty"`
}

// DeleteIndicatorsRequest is the request body for indicators/delete. The API
// has no by-ID delete: callers must construct a filter that matches the
// indicator(s) to remove (typically `{field:"indicator",operator:"EQ",value:<id>}`).
type DeleteIndicatorsRequest struct {
	Filters []IndicatorFilter `json:"filters"`
}

// InsertIndicatorsResponse is returned by /indicators/insert. Created
// records appear in AddedObjects; overwritten records appear in
// UpdatedObjects. Per-record failures appear in Errors keyed by the
// original index inside the batch — Errors items are objects, not strings,
// despite what the OpenAPI spec for /insert claims.
type InsertIndicatorsResponse struct {
	AddedObjects   []IndicatorInsertResult `json:"added_objects,omitempty"`
	UpdatedObjects []IndicatorInsertResult `json:"updated_objects,omitempty"`
	Errors         []IndicatorInsertError  `json:"errors,omitempty"`
}

// IndicatorInsertResult is one entry inside AddedObjects/UpdatedObjects
// from /indicators/insert. ID is the server-assigned rule_id; Status is a
// human-readable message describing what was done.
type IndicatorInsertResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// IndicatorInsertError is one entry inside Errors from /indicators/insert.
// Index is the zero-based position of the failing record inside the
// original request batch; Status is the server-supplied failure reason.
type IndicatorInsertError struct {
	Index  int    `json:"index"`
	Status string `json:"status"`
}

// ListIndicatorsResponse is returned by /indicators/get.
type ListIndicatorsResponse struct {
	ObjectsCount int         `json:"objects_count"`
	ObjectsType  string      `json:"objects_type"`
	Objects      []Indicator `json:"objects"`
}

// DeleteIndicatorsResponse is returned by /indicators/delete. Objects
// carries the server-assigned rule_ids that were removed; ObjectsCount is
// the length of that slice. Delete is idempotent — deleting a non-existent
// indicator returns `{objects_count: 0, objects: []}` with no error.
type DeleteIndicatorsResponse struct {
	ObjectsCount int   `json:"objects_count"`
	Objects      []int `json:"objects"`
}
