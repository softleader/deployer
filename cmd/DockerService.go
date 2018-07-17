package cmd

import (
	"strings"
	"strconv"
	"github.com/softleader/deployer/models"
)

type DockerService struct {
}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func (ds *DockerService) Ls() (arg string, out string, err error) {
	return Exec(&Options{}, "docker service ls", "--format '{{.Replicas}}'")
}

func (ds *DockerService) Update(service string, args ...string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service update", service, "-d", "--force", strings.Join(args, " "))
}

func (ds *DockerService) Inspect(service string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service inspect", service)
}

func (ds *DockerService) Logs(service string, tail int) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service logs --tail", strconv.Itoa(tail), service)
}

func (ds *DockerService) Rm(service string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service rm", service)
}

func (ds *DockerService) Ps(id string) (s []models.DockerServicePs, err error) {
	_, out, err := ps(id)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerServicePs(line))
		}
	}
	return
}

func ps(id string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service ps", id, "--no-trunc", "--format '{{.ID}};{{.Name}};{{.Image}};{{.Node}};{{.DesiredState}};{{.CurrentState}};{{.Error}}'")
}

func (ds *DockerService) GetCreatedTimeOfFirstServiceInStack(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker service inspect $(docker stack services", stack, "--format '{{.ID}}' | sed -n 1p) --format {{.CreatedAt}}")
}
