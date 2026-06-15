// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package enums

// ==============================================================================
// Indicator (IOC) Enums
// ==============================================================================
//
// Mirrors the enum sets declared by the `/public_api/v1/indicators/insert`
// schema. Underlying Go type is `string`, so callers can hold any value the
// API surfaces — including ones the OpenAPI spec doesn't yet document — but
// the named constants are the documented set and should be used as the
// default. Slice helpers below return the canonical set for use as switch
// cases, validator inputs, etc.

// ----------------------------------------------------------------------------
// Indicator Type
// ----------------------------------------------------------------------------

// IndicatorType is the kind of value an indicator describes.
//
// Note: URL is omitted from the OpenAPI insert enum but is present in the
// live API — verified by GET responses on a live tenant returning
// `"type": "URL"` records.
type IndicatorType string

const (
	IndicatorTypeHash       IndicatorType = "HASH"
	IndicatorTypeIP         IndicatorType = "IP"
	IndicatorTypePath       IndicatorType = "PATH"
	IndicatorTypeDomainName IndicatorType = "DOMAIN_NAME"
	IndicatorTypeFilename   IndicatorType = "FILENAME"
	IndicatorTypeMixed      IndicatorType = "MIXED"
	IndicatorTypeURL        IndicatorType = "URL"
)

// IndicatorTypes returns the IndicatorType values accepted by the live API.
// URL is included despite being missing from the OpenAPI spec.
func IndicatorTypes() []IndicatorType {
	return []IndicatorType{
		IndicatorTypeHash,
		IndicatorTypeIP,
		IndicatorTypePath,
		IndicatorTypeDomainName,
		IndicatorTypeFilename,
		IndicatorTypeMixed,
		IndicatorTypeURL,
	}
}

// ----------------------------------------------------------------------------
// Indicator Severity
// ----------------------------------------------------------------------------

// IndicatorSeverity is the severity rating assigned to an indicator. The
// namespaced `SEV_NNN_X` encoding is what the live API accepts on both
// reads and writes.
//
// Note: SEV_050_CRITICAL is omitted from the OpenAPI severity enum for
// /indicators/insert but is present in the Cortex UI and accepted by the
// live API.
type IndicatorSeverity string

const (
	IndicatorSeverityInfo     IndicatorSeverity = "SEV_010_INFO"
	IndicatorSeverityLow      IndicatorSeverity = "SEV_020_LOW"
	IndicatorSeverityMedium   IndicatorSeverity = "SEV_030_MEDIUM"
	IndicatorSeverityHigh     IndicatorSeverity = "SEV_040_HIGH"
	IndicatorSeverityCritical IndicatorSeverity = "SEV_050_CRITICAL"
)

// IndicatorSeverities returns the IndicatorSeverity values
func IndicatorSeverities() []IndicatorSeverity {
	return []IndicatorSeverity{
		IndicatorSeverityInfo,
		IndicatorSeverityLow,
		IndicatorSeverityMedium,
		IndicatorSeverityHigh,
		IndicatorSeverityCritical,
	}
}

// ----------------------------------------------------------------------------
// Indicator Reputation
// ----------------------------------------------------------------------------

// IndicatorReputation is a categorical assessment of the indicator's
// reputation.
type IndicatorReputation string

const (
	IndicatorReputationGood         IndicatorReputation = "GOOD"
	IndicatorReputationBad          IndicatorReputation = "BAD"
	IndicatorReputationSuspicious   IndicatorReputation = "SUSPICIOUS"
	IndicatorReputationUnknown      IndicatorReputation = "UNKNOWN"
	IndicatorReputationNoReputation IndicatorReputation = "NO_REPUTATION"
)

// IndicatorReputations returns the documented IndicatorReputation values in
// declaration order.
func IndicatorReputations() []IndicatorReputation {
	return []IndicatorReputation{
		IndicatorReputationGood,
		IndicatorReputationBad,
		IndicatorReputationSuspicious,
		IndicatorReputationUnknown,
		IndicatorReputationNoReputation,
	}
}

// ----------------------------------------------------------------------------
// Indicator Reliability
// ----------------------------------------------------------------------------

// IndicatorReliability is a single-character reliability rating where A is
// the most reliable and G the least. The API also accepts the empty string
// to mean "unset".
type IndicatorReliability string

const (
	IndicatorReliabilityA IndicatorReliability = "A"
	IndicatorReliabilityB IndicatorReliability = "B"
	IndicatorReliabilityC IndicatorReliability = "C"
	IndicatorReliabilityD IndicatorReliability = "D"
	IndicatorReliabilityE IndicatorReliability = "E"
	IndicatorReliabilityF IndicatorReliability = "F"
	IndicatorReliabilityG IndicatorReliability = "G"
)

// IndicatorReliabilities returns the documented IndicatorReliability values
// in declaration order.
func IndicatorReliabilities() []IndicatorReliability {
	return []IndicatorReliability{
		IndicatorReliabilityA,
		IndicatorReliabilityB,
		IndicatorReliabilityC,
		IndicatorReliabilityD,
		IndicatorReliabilityE,
		IndicatorReliabilityF,
		IndicatorReliabilityG,
	}
}
