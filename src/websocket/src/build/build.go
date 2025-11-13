package build

import (
	"fmt"
	"os"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/runFlags"
	"motorv2/src/websocket/src/runFlags"
)

var Version string

func Init() {
	runFlags.Parse()

	if *runFlags.Flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}
}
