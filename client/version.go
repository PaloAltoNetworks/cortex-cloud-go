// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"runtime"
)

var (
	GitCommit = "NOCOMMIT"
	GoVersion = runtime.Version()
	BuildDate = ""
)
