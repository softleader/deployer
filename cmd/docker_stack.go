package cmd

import (
	"strconv"
	"strings"
	"github.com/softleader/deployer/models"
)

type DockerStack struct {
	login string
}

func NewDockerStack(login string) *DockerStack {
	return &DockerStack{login: login}
}

func (ds *DockerStack) Ls() (s []models.DockerStackLs, err error) {
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
	return Exec(&Options{}, "docker stack ls --format '{{.Name}};{{.Services}}'")
}

func (ds *DockerStack) Services(name string) (s []models.DockerStackServices, err error) {
	_, out, err := services(name)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerStackServices(line))
		}
	}
	return
}

func services(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker stack services", stack, "--format '{{.ID}};{{.Name}};{{.Mode}};{{.Replicas}};{{.Image}};{{.Ports}}'")
}

func (ds *DockerStack) RmLike(stack string) (arg string, out string, err error) {
	return Exec(&Options{}, "docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func deploy(opts *Options, stack string, file string, login string) (arg string, out string, err error) {
	return Exec(opts, login, "docker stack deploy -c", file, stack, "--with-registry-auth")
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
		deploy(opts, strings.Join(stack, "-"), y.Path, ds.login)
	}
	return nil
}
