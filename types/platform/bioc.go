// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"encoding/json"

	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
)

// ----------------------------------------------------------------------------
// BIOC
// ----------------------------------------------------------------------------

// BIOC describes a single Behavioral Indicator of Compromise as accepted by
// /bioc/insert and returned by /bioc/get.
//
// The first ten fields are read+write. The remaining four are read-only:
// they appear in GET responses but are server-assigned and not part of the
// insert schema. CreationTime/ModificationTime/Source/NumberOfIssues are
// undocumented in the OpenAPI spec for /bioc/get but populated by the live
// API — observed against a live tenant.
//
// The Indicator field is polymorphic: when IsXQL is true it carries a raw
// XQL query string; otherwise it carries a structured filter-AST object
// (runOnCGO / investigationType / investigation.{TYPE}.filter.{AND|OR}[]).
// json.RawMessage lets the SDK pass both shapes through without imposing a
// (necessarily incomplete) Go struct for the AST.
type BIOC struct {
	// --- Read + Write fields ---
	RuleID                  int                `json:"rule_id,omitempty"`
	Name                    string             `json:"name"`
	Type                    enums.BIOCType     `json:"type"`
	Severity                enums.BIOCSeverity `json:"severity"`
	Status                  enums.BIOCStatus   `json:"status"`
	Comment                 string             `json:"comment"`
	IsXQL                   bool               `json:"is_xql"`
	Indicator               json.RawMessage    `json:"indicator"`
	MitreTacticIDAndName    []string           `json:"mitre_tactic_id_and_name"`
	MitreTechniqueIDAndName []string           `json:"mitre_technique_id_and_name"`

	// --- Read-only fields (server-assigned, returned by /bioc/get) ---
	CreationTime     int64  `json:"creation_time,omitempty"`
	ModificationTime int64  `json:"modification_time,omitempty"`
	Source           string `json:"source,omitempty"`
	NumberOfIssues   int    `json:"number_of_issues,omitempty"`
}

// BIOCFilter is a single filter expression used by bioc/get and bioc/delete.
// Filter `value` is polymorphic — booleans (`is_xql`) must be JSON booleans,
// integers (`rule_id`) must be JSON numbers, everything else passes through
// as a string.
//
// Note: the OpenAPI filter `field` enum lists name/severity/type/is_xql/
// comment/status/indicator/mitre_*; rule_id is undocumented but accepted on
// EQ by both bioc/get and bioc/delete (verified against a live tenant).
// FindBIOCByID and DeleteBIOCsByRuleID rely on this behavior.
type BIOCFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// ListBIOCsRequest is the request body for bioc/get. ExtendedView toggles
// inclusion of additional fields in the response (in practice the live API
// returns the same field set either way for BIOCs created via /insert).
type ListBIOCsRequest struct {
	Filters      []BIOCFilter `json:"filters,omitempty"`
	ExtendedView bool         `json:"extended_view,omitempty"`
	SearchFrom   int          `json:"search_from,omitempty"`
	SearchTo     int          `json:"search_to,omitempty"`
}

// DeleteBIOCsRequest is the request body for bioc/delete. The API has no
// by-ID delete: callers must construct a filter that matches the BIOC(s) to
// remove. The undocumented `rule_id` filter field is the only safe identity
// key — BIOC names are not unique on a tenant.
type DeleteBIOCsRequest struct {
	Filters []BIOCFilter `json:"filters"`
}

// InsertBIOCsResponse is returned by /bioc/insert. Created records appear in
// AddedObjects; overwritten records appear in UpdatedObjects. Per-record
// failures appear in Errors keyed by the original index inside the batch.
//
// Unlike /indicators/insert, the BIOC endpoint returns HTTP 400 (not 200)
// when any record fails validation, but the body still uses the success
// shape. The SDK insert helper handles this and returns the typed response
// rather than the raw HTTP error.
type InsertBIOCsResponse struct {
	AddedObjects   []BIOCInsertResult `json:"added_objects,omitempty"`
	UpdatedObjects []BIOCInsertResult `json:"updated_objects,omitempty"`
	Errors         []BIOCInsertError  `json:"errors,omitempty"`
}

// BIOCInsertResult is one entry inside AddedObjects/UpdatedObjects from
// /bioc/insert.
type BIOCInsertResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// BIOCInsertError is one entry inside Errors from /bioc/insert. Index is the
// zero-based position of the failing record inside the original request
// batch; Status is the server-supplied failure reason. The OpenAPI spec
// declares Errors items as strings, but the live API returns objects with
// {index, status} (verified against a live tenant).
type BIOCInsertError struct {
	Index  int    `json:"index"`
	Status string `json:"status"`
}

// ListBIOCsResponse is returned by /bioc/get.
type ListBIOCsResponse struct {
	ObjectsCount int    `json:"objects_count"`
	ObjectsType  string `json:"objects_type"`
	Objects      []BIOC `json:"objects"`
}

// DeleteBIOCsResponse is returned by /bioc/delete. Objects carries the
// server-assigned rule_ids that were removed; ObjectsCount is the length of
// that slice. Delete is idempotent — deleting a non-existent BIOC returns
// `{objects_count: 0, objects: []}` with no error.
type DeleteBIOCsResponse struct {
	ObjectsCount int   `json:"objects_count"`
	Objects      []int `json:"objects"`
}
