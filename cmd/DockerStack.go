package cmd

import (
	"strconv"
	"strings"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/app"
)

type DockerStack struct {
	app.Registry
}

func NewDockerStack(registry app.Registry) *DockerStack {
	return &DockerStack{Registry: registry}
}

func (ds *DockerStack) Ls() ([][]string, error) {
	_, out, err := ls()
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, strings.Split(line, ";"))
		}
	}
	return s, err
}

func ls() (arg string, out string, err error) {
	return Exec(&Options{}, "docker stack ls --format '{{.Name}};{{.Services}}'")
}

func (ds *DockerStack) Services(name string) ([][]string, error) {
	_, out, err := services(name)
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, strings.Split(line, ";"))
		}
	}
	return s, err
}

func services(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker stack services", stack, "--format '{{.ID}};{{.Name}};{{.Mode}};{{.Replicas}};{{.Image}};{{.Ports}}'")
}

func (ds *DockerStack) RmLike(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func deploy(opts *Options, stack string, file string, registry app.Registry) (arg string, out string, err error) {
	return Exec(opts, registry.Login(), "docker stack deploy -c", file, stack, "--with-registry-auth")
}

func (ds *DockerStack) Deploy(opts *Options, yamls []models.Yaml, d *models.Deploy) error {
	for _, y := range yamls {
		stack := []string{d.Project}
		if d.Dev.IpAddress != "" {
			stack = append(stack, strconv.Itoa(d.Dev.Port))
		}
		if y.Group != "" {
			stack = append(stack, y.Group)
		}
		deploy(opts, strings.Join(stack, "-"), y.Path, ds.Registry)
	}
	return nil
}
