package cmd

type DockerStack struct {
	sh Shell
}

func NewDockerStack(sh Shell) *DockerStack {
	return &DockerStack{sh: sh}
}

func (ds *DockerStack) Ls() (string, string, error) {
	return ds.sh.Exec(&Options{}, "docker stack ls --format '{{.Name}};{{.Services}}'")
}

func (ds *DockerStack) Services(stack string) (string, string, error) {
	return ds.sh.Exec(&Options{}, "docker stack services", stack, "--format '{{.ID}};{{.Name}};{{.Mode}};{{.Replicas}};{{.Image}};{{.Ports}}'")
}

func (ds *DockerStack) Ps(id string) (string, string, error) {
	return ds.sh.Exec(&Options{}, "docker service ps", id, "--no-trunc", "--format '{{.ID}};{{.Name}};{{.Image}};{{.Node}};{{.DesiredState}};{{.CurrentState}};{{.Error}}'")
}

//func (ds *DockerStack) RmStack(stack string) (string, string, error) {
//	return ds.sh.Exec(&Options{}, "docker stack rm", stack)
//}

func (ds *DockerStack) RmStackLike(stack string) (string, string, error) {
	return ds.sh.Exec(&Options{}, "docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func (ds *DockerStack) RmService(service string) (string, string, error) {
	return ds.sh.Exec(&Options{}, "docker service rm", service)
}

func (ds *DockerStack) Deploy(opts *Options, stack string, file string) (string, string, error) {
	return ds.sh.Exec(opts, "docker stack deploy -c", file, stack)
}
