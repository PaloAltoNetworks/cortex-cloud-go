// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"runtime"
)

var (
	Version             = ""
	GitCommit           = "NOCOMMIT"
	CortexServerVersion = ""
	CortexPAPIVersion   = ""
	GoVersion           = runtime.Version()
	BuildDate           = ""
)
