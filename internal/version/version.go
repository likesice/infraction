package version

import "fmt"

var version string = "1.0.0-alpha" // set by linkerflags at build-time

func GetVersionString() string {
	return fmt.Sprintf("version %v", version)
}

func GetVersion() string {
	return version
}
