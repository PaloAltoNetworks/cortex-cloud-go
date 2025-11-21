// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

const (
	// Version is the semantic version of this module
	Version = "1.0.0"

	// ModuleName is the canonical name of this module
	ModuleName = "cwp"
)

var (
	// GitCommit is the git commit hash (set via ldflags)
	GitCommit = "dev"

	// CortexServerVersion is the target Cortex server version
	CortexServerVersion = "unknown"

	// CortexPAPIVersion is the target PAPI version
	CortexPAPIVersion = "unknown"

	// GoVersion is the Go version used to build
	GoVersion = "unknown"

	// BuildDate is the build timestamp
	BuildDate = "unknown"
)
