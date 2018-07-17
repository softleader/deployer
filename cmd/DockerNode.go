package cmd

import (
	"github.com/softleader/deployer/models"
	"strings"
)

type DockerNode struct {
}

func NewDockerNode() *DockerNode {
	return &DockerNode{}
}

func (ds *DockerNode) Ls() (s []models.DockerNodeLs, err error) {
	_, out, err := nodeLs()
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerNodeLs(line))
		}
	}
	return
}

func nodeLs() (arg string, out string, err error) {
	return Exec(&Options{}, "docker node ls", "--format '{{.Hostname}};{{.Status}};{{.Availability}}'")
}
