package cmd

type DockerService struct {
	sh Shell
}

func NewDockerService(sh Shell) *DockerService {
	return &DockerService{sh: sh}
}

func (ds *DockerService) Rm(service string) (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker service rm", service)
}

func (ds *DockerService) Ps(id string) (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker service ps", id, "--no-trunc", "--format '{{.ID}};{{.Name}};{{.Image}};{{.Node}};{{.DesiredState}};{{.CurrentState}};{{.Error}}'")
}

func (ds *DockerService) GetCreatedTimeOfFirstServiceInStack(stack string) (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker service inspect $(docker stack services", stack, "--format '{{.ID}}' | sed -n 1p) --format {{.CreatedAt}}")
}
