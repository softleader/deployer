package cmd

type DockerStack struct {
	sh Shell
}

func NewDockerStack(sh Shell) *DockerStack {
	return &DockerStack{sh: sh}
}

func (ds *DockerStack) Ls() (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker stack ls --format '{{.Name}};{{.Services}}'")
}

func (ds *DockerStack) Services(stack string) (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker stack services", stack, "--format '{{.ID}};{{.Name}};{{.Mode}};{{.Replicas}};{{.Image}};{{.Ports}}'")
}

//func (ds *DockerStack) RmStack(stack string) (arg string, out string, err error) {
//	return ds.sh.Exec(&Options{}, "docker stack rm", stack)
//}

func (ds *DockerStack) RmLike(stack string) (arg string, out string, err error) {
	return ds.sh.Exec(&Options{}, "docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func (ds *DockerStack) Deploy(opts *Options, stack string, file string) (arg string, out string, err error) {
	return ds.sh.Exec(opts, "docker stack deploy -c", file, stack)
}
