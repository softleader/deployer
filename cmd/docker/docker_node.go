package docker

import (
	"github.com/softleader/deployer/models"
	"strings"
	"github.com/softleader/deployer/cmd"
)

func NodeLs() (s []models.DockerNodeLs, err error) {
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
	return cmd.Exec(&cmd.Options{}, "docker node ls", "--format '{{.Hostname}};{{.Status}};{{.Availability}}'")
}
