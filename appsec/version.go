package appsec

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	GitCommit 	= "NOCOMMIT"
	GoVersion 	= runtime.Version()
	BuildDate   = ""
)

func initVersion() {  
    info, ok := debug.ReadBuildInfo()  
	fmt.Printf("%+v", info)
    if !ok {  
       return  
    }  
    modified := false
    for _, setting := range info.Settings {  
       switch setting.Key {  
       case "vcs.revision":  
          GitCommit = setting.Value  
       case "vcs.time":  
          BuildDate = setting.Value  
       case "vcs.modified":  
          modified = true  
       }  
    }  
    if modified {  
       GitCommit += "+CHANGES"
    }  
}
