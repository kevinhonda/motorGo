package build

import (
	"fmt"
	"os"
	"motorv2/src/runFlags"
)

var Version string

func Init() {
	runFlags.Parse()

	if *runFlags.Flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}
}
