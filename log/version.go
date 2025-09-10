package log

import (
	"runtime"
)

var (
	GitCommit = "NOCOMMIT"
	GoVersion = runtime.Version()
	BuildDate = ""
)
