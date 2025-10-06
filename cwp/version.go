// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cwp

import (
	"runtime"
)

var (
	GitCommit           = "NOCOMMIT"
	CortexServerVersion = ""
	CortexPAPIVersion   = ""
	GoVersion           = runtime.Version()
	BuildDate           = ""
)
