package docker

import (
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/models"
	"strconv"
	"strings"
)

func StackLs() (s []models.DockerStackLs, err error) {
	_, out, err := stackLs()
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerStackLs(line))
		}
	}
	return
}

func stackLs() (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker stack ls --format '{{.Name}};{{.Services}}'")
}

func StackServices(name string) (s []models.DockerStackServices, err error) {
	_, out, err := stackServices(name)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			ss := models.NewDockerStackServices(line)
			if _, spec, err := ServiceSpec(ss.Name); err == nil {
				ss.Labels = spec.Labels
			}
			s = append(s, ss)
		}
	}

	return
}

func stackServices(stack string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker stack services", stack, "--format '{{.ID}};{{.Name}};{{.Mode}};{{.Replicas}};{{.Image}};{{.Ports}}'")
}

func StackRmLike(stack string) (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, "docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func stackDeploy(opts *cmd.Options, stack string, file string, login string) (arg string, out string, err error) {
	return cmd.Exec(opts, login, "docker stack deploy -c", file, stack, "--with-registry-auth")
}

func StackDeploy(opts *cmd.Options, yamls []models.Yaml, d *models.Deploy, login string) error {
	for _, y := range yamls {
		stack := []string{d.Project}
		if d.Dev.IpAddress != "" {
			stack = append(stack, strconv.Itoa(d.Dev.Port))
		}
		if y.Group != "" {
			stack = append(stack, y.Group)
		}
		stackDeploy(opts, strings.Join(stack, "-"), y.Path, login)
	}
	return nil
}
