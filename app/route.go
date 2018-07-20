package app

import "github.com/softleader/deployer/cmd"

// it's just a collection of all commands for each route use
type Route struct {
	*Args
	*Workspace
	cmd.GenYaml
	cmd.Gpm
}
