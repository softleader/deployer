package docker

import (
	"encoding/json"
	"fmt"
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/models"
	"strconv"
	"strings"
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

type Service struct {
	Spec Spec
}

type Spec struct {
	Name   string
	Labels map[string]string `json:"Labels"`
}

func ServiceSpec(service string) (string, *Spec, error) {
	arg, out, err := cmd.Exec(&cmd.Options{}, "docker service inspect", service, "-f", "'{{json .}}'")
	if err != nil {
		return arg, nil, err
	}
	svc := &Service{}
	if err := json.Unmarshal([]byte(out), svc); err != nil {
		return arg, nil, err
	}
	return arg, &svc.Spec, err
}

func ServiceFilter(params map[string]string) (arg string, out string, err error) {
	args := []string{"docker service ls --format '{{json .}}'"}
	for key, val := range params {
		args = append(args, "-f", fmt.Sprintf("%s=%s", key, val))
	}
	return cmd.Exec(&cmd.Options{}, args...)
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
