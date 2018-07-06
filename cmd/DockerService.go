package cmd

import (
	"strings"
)

type DockerService struct {
}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func (ds *DockerService) Inspect(service string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service inspect", service)
}

func (ds *DockerService) Rm(service string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service rm", service)
}

func (ds *DockerService) Ps(id string) ([][]string, error) {
	_, out, err := ps(id)
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fields := strings.Split(line, ";")
			fields[2] = strings.Split(fields[2], "@sha256")[0]
			s = append(s, fields)
		}
	}
	return s, err
}

func ps(id string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service ps", id, "--no-trunc", "--format '{{.ID}};{{.Name}};{{.Image}};{{.Node}};{{.DesiredState}};{{.CurrentState}};{{.Error}}'")
}

func (ds *DockerService) GetCreatedTimeOfFirstServiceInStack(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service inspect $(docker stack services", stack, "--format '{{.ID}}' | sed -n 1p) --format {{.CreatedAt}}")
}
