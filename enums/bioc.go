// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package enums

// ==============================================================================
// BIOC Enums
// ==============================================================================
//
// Mirrors the enum sets declared by the `/public_api/v1/bioc/insert` schema.
// Underlying Go type is `string`, so callers can hold any value the API
// surfaces — including ones the OpenAPI spec doesn't yet document — but the
// named constants are the documented set and should be used as the default.

// ----------------------------------------------------------------------------
// BIOC Type
// ----------------------------------------------------------------------------

// BIOCType is the category a BIOC rule falls under. The 16 values mirror the
// `/public_api/v1/bioc/insert` request schema and match the casing the live
// API returns on GET.
type BIOCType string

const (
	BIOCTypeOther                    BIOCType = "OTHER"
	BIOCTypePersistence              BIOCType = "PERSISTENCE"
	BIOCTypeEvasion                  BIOCType = "EVASION"
	BIOCTypeTampering                BIOCType = "TAMPERING"
	BIOCTypeFileTypeObfuscation      BIOCType = "FILE_TYPE_OBFUSCATION"
	BIOCTypePrivilegeEscalation      BIOCType = "PRIVILEGE_ESCALATION"
	BIOCTypeCredentialAccess         BIOCType = "CREDENTIAL_ACCESS"
	BIOCTypeLateralMovement          BIOCType = "LATERAL_MOVEMENT"
	BIOCTypeExecution                BIOCType = "EXECUTION"
	BIOCTypeCollection               BIOCType = "COLLECTION"
	BIOCTypeExfiltration             BIOCType = "EXFILTRATION"
	BIOCTypeInfiltration             BIOCType = "INFILTRATION"
	BIOCTypeDropper                  BIOCType = "DROPPER"
	BIOCTypeFilePrivilegeManipulation BIOCType = "FILE_PRIVILEGE_MANIPULATION"
	BIOCTypeReconnaissance           BIOCType = "RECONNAISSANCE"
	BIOCTypeDiscovery                BIOCType = "DISCOVERY"
)

// BIOCTypes returns the documented BIOCType values in declaration order.
func BIOCTypes() []BIOCType {
	return []BIOCType{
		BIOCTypeOther,
		BIOCTypePersistence,
		BIOCTypeEvasion,
		BIOCTypeTampering,
		BIOCTypeFileTypeObfuscation,
		BIOCTypePrivilegeEscalation,
		BIOCTypeCredentialAccess,
		BIOCTypeLateralMovement,
		BIOCTypeExecution,
		BIOCTypeCollection,
		BIOCTypeExfiltration,
		BIOCTypeInfiltration,
		BIOCTypeDropper,
		BIOCTypeFilePrivilegeManipulation,
		BIOCTypeReconnaissance,
		BIOCTypeDiscovery,
	}
}

// ----------------------------------------------------------------------------
// BIOC Severity
// ----------------------------------------------------------------------------

// BIOCSeverity is the severity rating assigned to a BIOC. The namespaced
// `SEV_NNN_X` encoding matches the live API.
//
// Note: SEV_050_CRITICAL is omitted from the OpenAPI severity enum for
// /bioc/insert but is present in the Cortex UI and accepted by the live
// API — verified by inserting a record with severity=SEV_050_CRITICAL on
// a live tenant and round-tripping it via /bioc/get. Same behavior as
// IndicatorSeverity.
type BIOCSeverity string

const (
	BIOCSeverityInfo     BIOCSeverity = "SEV_010_INFO"
	BIOCSeverityLow      BIOCSeverity = "SEV_020_LOW"
	BIOCSeverityMedium   BIOCSeverity = "SEV_030_MEDIUM"
	BIOCSeverityHigh     BIOCSeverity = "SEV_040_HIGH"
	BIOCSeverityCritical BIOCSeverity = "SEV_050_CRITICAL"
)

// BIOCSeverities returns the BIOCSeverity values accepted by the live API.
// SEV_050_CRITICAL is included despite being missing from the OpenAPI spec.
func BIOCSeverities() []BIOCSeverity {
	return []BIOCSeverity{
		BIOCSeverityInfo,
		BIOCSeverityLow,
		BIOCSeverityMedium,
		BIOCSeverityHigh,
		BIOCSeverityCritical,
	}
}

// ----------------------------------------------------------------------------
// BIOC Status
// ----------------------------------------------------------------------------

// BIOCStatus toggles whether a BIOC rule is active. The OpenAPI schema
// declares the enum in lowercase (`enabled`/`disabled`), and the live API
// preserves whatever casing was written — records created via the UI come
// back uppercase, records inserted via the SDK in lowercase come back
// lowercase. Callers comparing status should normalize first; the SDK
// constants are the lowercase canonical form.
type BIOCStatus string

const (
	BIOCStatusEnabled  BIOCStatus = "enabled"
	BIOCStatusDisabled BIOCStatus = "disabled"
)

// BIOCStatuses returns the documented BIOCStatus values in declaration
// order.
func BIOCStatuses() []BIOCStatus {
	return []BIOCStatus{
		BIOCStatusEnabled,
		BIOCStatusDisabled,
	}
}
