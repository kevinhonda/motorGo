package msServices

import //"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/cmd"
"motorv2/src/websocket/pkg/cmd"

func GetStarted() ([]string, error) {
	return cmd.RunOutput("wmic", "service", "where", "started=true", "get", "name")
}

func GetStopped() ([]string, error) {
	return cmd.RunOutput("wmic", "service", "where", "started=false", "get", "name")
}

func GetLike(text string) ([]string, error) {
	return cmd.RunOutput("wmic", "service", "where", "name like \"%"+text+"%\"", "get", "name,started")
}

func Start(name string) error {
	return cmd.Run("net", "start", name)
}

func Stop(name string) error {
	return cmd.Run("net", "stop", name)
}

func Restart(name string) error {
	err := Start(name)
	if err != nil {
		return err
	}

	return Stop(name)
}
