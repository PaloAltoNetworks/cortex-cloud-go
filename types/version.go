// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package types contains shared data structures used across different API
// modules.
package types

// TODO: move all types that are use in only a single package to those packages, then move all files to the root of the "types" directory so there are no nested directories or non-shared types.

var (
	GitCommit           = "NOCOMMIT"
	CortexServerVersion = "UNKNOWN"
	CortexPAPIVersion   = "UNKNOWN"
	GoVersion           = "UNKNOWN"
	BuildDate           = "UNKNOWN"
)
