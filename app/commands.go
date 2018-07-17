package app

import "github.com/softleader/deployer/cmd"

// it just a collection of all commands
type Commands struct {
	cmd.DockerNode
	cmd.DockerStack
	cmd.DockerService
	cmd.DockerStats
	cmd.GenYaml
	cmd.Gpm
}
