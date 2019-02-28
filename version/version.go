package version

import (
	"fmt"
	"runtime"
)

// GitCommit The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version The main version number that is being run at the moment.
var Version = ""

// BuildDate is the date the binary was built
var BuildDate = ""

// GoVersion is the version of go use to build the binary
var GoVersion = runtime.Version()

// OsArch is the os and processor architecture the binary was built for
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
