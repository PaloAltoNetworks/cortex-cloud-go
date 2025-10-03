// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package errors

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
