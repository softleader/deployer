package docker

import (
	"strings"
	"strconv"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/cmd"
)

func ServiceLs() (s []models.DockerServiceLs, err error) {
	_, out, err := serviceLs()
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerServiceLs(line))
		}
	}
	return
}

func serviceLs() (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service ls", "--format '{{.Replicas}}'")
}

func ServiceUpdate(service string, args ...string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service update", service, "-d", "--force", strings.Join(args, " "))
}

func ServiceInspect(service string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service inspect", service)
}

func ServiceLogs(service string, tail int) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service logs --tail", strconv.Itoa(tail), service)
}

func ServiceRm(service string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service rm", service)
}

func ServicePs(id string) (s []models.DockerServicePs, err error) {
	_, out, err := servicePs(id)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerServicePs(line))
		}
	}
	return
}

func servicePs(id string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service ps", id, "--no-trunc", "--format '{{.ID}};{{.Name}};{{.Image}};{{.Node}};{{.DesiredState}};{{.CurrentState}};{{.Error}}'")
}

func ServiceGetCreatedTimeOfFirstServiceInStack(stack string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker service inspect $(docker stack services", stack, "--format '{{.ID}}' | sed -n 1p) --format {{.CreatedAt}}")
}
