// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
	"testing"
)

// TestErrorExtraField_UnmarshalJSON_StringFormat tests unmarshaling when err_extra is a string
func TestErrorExtraField_UnmarshalJSON_StringFormat(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		wantLen     int
		wantMessage string
		wantType    string
	}{
		{
			name:        "non-empty string",
			jsonData:    `{"err_code": 400, "err_msg": "Bad Request", "err_extra": "Missing at least one required parameter: ['evaluation_frequency']."}`,
			wantLen:     1,
			wantMessage: "Missing at least one required parameter: ['evaluation_frequency'].",
			wantType:    "string_error",
		},
		{
			name:     "empty string",
			jsonData: `{"err_code": 400, "err_msg": "Bad Request", "err_extra": ""}`,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reply CortexCloudAPIErrorReply
			err := json.Unmarshal([]byte(tt.jsonData), &reply)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			values := reply.Extra.Values()
			if len(values) != tt.wantLen {
				t.Errorf("Expected %d error extra values, got %d", tt.wantLen, len(values))
			}

			if tt.wantLen > 0 {
				if values[0].Message != tt.wantMessage {
					t.Errorf("Expected message %q, got %q", tt.wantMessage, values[0].Message)
				}
				if values[0].Type != tt.wantType {
					t.Errorf("Expected type %q, got %q", tt.wantType, values[0].Type)
				}
			}
		})
	}
}

// TestErrorExtraField_UnmarshalJSON_ArrayFormat tests unmarshaling when err_extra is an array
func TestErrorExtraField_UnmarshalJSON_ArrayFormat(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantLen  int
	}{
		{
			name: "array with single object",
			jsonData: `{
				"err_code": 400,
				"err_msg": "Validation Error",
				"err_extra": [
					{
						"type": "missing",
						"loc": ["body", "name"],
						"msg": "Field required",
						"input": null
					}
				]
			}`,
			wantLen: 1,
		},
		{
			name: "array with multiple objects",
			jsonData: `{
				"err_code": 400,
				"err_msg": "Validation Error",
				"err_extra": [
					{
						"type": "missing",
						"loc": ["body", "name"],
						"msg": "Field required",
						"input": null
					},
					{
						"type": "string_too_short",
						"loc": ["body", "description"],
						"msg": "String should have at least 1 character",
						"input": "",
						"ctx": {"min_length": 1}
					}
				]
			}`,
			wantLen: 2,
		},
		{
			name:     "empty array",
			jsonData: `{"err_code": 400, "err_msg": "Bad Request", "err_extra": []}`,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reply CortexCloudAPIErrorReply
			err := json.Unmarshal([]byte(tt.jsonData), &reply)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			values := reply.Extra.Values()
			if len(values) != tt.wantLen {
				t.Errorf("Expected %d error extra values, got %d", tt.wantLen, len(values))
			}
		})
	}
}

// TestErrorExtraField_UnmarshalJSON_ArrayFormatDetails tests array format with detailed validation
func TestErrorExtraField_UnmarshalJSON_ArrayFormatDetails(t *testing.T) {
	jsonData := `{
		"err_code": 400,
		"err_msg": "Validation Error",
		"err_extra": [
			{
				"type": "missing",
				"loc": ["body", "name"],
				"msg": "Field required",
				"input": null
			}
		]
	}`

	var reply CortexCloudAPIErrorReply
	err := json.Unmarshal([]byte(jsonData), &reply)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	values := reply.Extra.Values()
	if len(values) != 1 {
		t.Fatalf("Expected 1 error extra value, got %d", len(values))
	}

	extra := values[0]
	if extra.Type != "missing" {
		t.Errorf("Expected type 'missing', got %q", extra.Type)
	}
	if extra.Message != "Field required" {
		t.Errorf("Expected message 'Field required', got %q", extra.Message)
	}
	if len(extra.Location) != 2 {
		t.Errorf("Expected 2 location elements, got %d", len(extra.Location))
	}
}

// TestErrorExtraField_UnmarshalJSON_InvalidFormat tests error handling for invalid formats
func TestErrorExtraField_UnmarshalJSON_InvalidFormat(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "number format",
			jsonData: `{"err_code": 400, "err_msg": "Bad Request", "err_extra": 123}`,
		},
		{
			name:     "boolean format",
			jsonData: `{"err_code": 400, "err_msg": "Bad Request", "err_extra": true}`,
		},
		{
			name:     "object format",
			jsonData: `{"err_code": 400, "err_msg": "Bad Request", "err_extra": {"key": "value"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reply CortexCloudAPIErrorReply
			err := json.Unmarshal([]byte(tt.jsonData), &reply)
			if err == nil {
				t.Error("Expected unmarshal to fail for invalid format, but it succeeded")
			}
		})
	}
}

// TestErrorExtraField_MarshalJSON tests marshaling back to JSON
func TestErrorExtraField_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		field    ErrorExtraField
		wantJSON string
	}{
		{
			name: "array format",
			field: ErrorExtraField{
				values: []CortexCloudAPIErrorExtra{
					{
						Type:    "missing",
						Message: "Field required",
					},
				},
			},
			wantJSON: `[{"type":"missing","loc":null,"msg":"Field required","input":null,"ctx":{}}]`,
		},
		{
			name: "empty array",
			field: ErrorExtraField{
				values: []CortexCloudAPIErrorExtra{},
			},
			wantJSON: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.field)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if string(data) != tt.wantJSON {
				t.Errorf("Expected JSON %q, got %q", tt.wantJSON, string(data))
			}
		})
	}
}

// TestCortexCloudAPIError_Error_StringFormat tests Error() method with string err_extra
func TestCortexCloudAPIError_Error_StringFormat(t *testing.T) {
	jsonData := `{
		"reply": {
			"err_code": 400,
			"err_msg": "The request contains invalid or missing parameters.",
			"err_extra": "Missing at least one required parameter: ['evaluation_frequency']."
		}
	}`

	var apiErr CortexCloudAPIError
	err := json.Unmarshal([]byte(jsonData), &apiErr)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	errStr := apiErr.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}

	// Verify the error string contains key information
	if !contains(errStr, "Error Code: 400") {
		t.Error("Error string should contain error code")
	}
	if !contains(errStr, "The request contains invalid or missing parameters.") {
		t.Error("Error string should contain error message")
	}
	if !contains(errStr, "Missing at least one required parameter") {
		t.Error("Error string should contain err_extra message")
	}
}

// TestCortexCloudAPIError_Error_ArrayFormat tests Error() method with array err_extra
func TestCortexCloudAPIError_Error_ArrayFormat(t *testing.T) {
	jsonData := `{
		"reply": {
			"err_code": 400,
			"err_msg": "Validation Error",
			"err_extra": [
				{
					"type": "missing",
					"loc": ["body", "name"],
					"msg": "Field required",
					"input": null
				}
			]
		}
	}`

	var apiErr CortexCloudAPIError
	err := json.Unmarshal([]byte(jsonData), &apiErr)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	errStr := apiErr.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}

	// Verify the error string contains key information
	if !contains(errStr, "Error Code: 400") {
		t.Error("Error string should contain error code")
	}
	if !contains(errStr, "Validation Error") {
		t.Error("Error string should contain error message")
	}
	if !contains(errStr, "Field required") {
		t.Error("Error string should contain validation message")
	}
}

// TestErrorExtraField_NullValue tests handling of null err_extra
func TestErrorExtraField_NullValue(t *testing.T) {
	jsonData := `{"err_code": 400, "err_msg": "Bad Request", "err_extra": null}`

	var reply CortexCloudAPIErrorReply
	err := json.Unmarshal([]byte(jsonData), &reply)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	values := reply.Extra.Values()
	if values != nil && len(values) != 0 {
		t.Errorf("Expected nil or empty slice for null err_extra, got %d values", len(values))
	}
}

// TestBackwardCompatibility tests that existing code patterns still work
func TestBackwardCompatibility(t *testing.T) {
	// Test that the Values() method returns a slice that can be ranged over
	jsonData := `{
		"err_code": 400,
		"err_msg": "Validation Error",
		"err_extra": [
			{
				"type": "missing",
				"loc": ["body", "name"],
				"msg": "Field required",
				"input": null
			}
		]
	}`

	var reply CortexCloudAPIErrorReply
	err := json.Unmarshal([]byte(jsonData), &reply)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// This is the pattern used in the Error() method
	count := 0
	for _, extra := range reply.Extra.Values() {
		count++
		if extra.Type == "" {
			t.Error("Expected non-empty type")
		}
	}

	if count != 1 {
		t.Errorf("Expected to iterate over 1 item, got %d", count)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
