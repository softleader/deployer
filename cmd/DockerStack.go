package cmd

type DockerStack struct {
	sh Sh
}

func NewDockerStack(sh Sh) DockerStack {
	return DockerStack{sh: sh}
}

func (ds DockerStack) Ls() (string, string, error) {
	return ds.sh.Exec("docker stack ls")
}

func (ds DockerStack) Services(stack string) (string, string, error) {
	return ds.sh.Exec("docker stack services", stack)
}

func (ds DockerStack) Rm(stack string) (string, string, error) {
	return ds.sh.Exec("docker stack rm", stack)
}

func (ds DockerStack) Deploy(stack string, file string) (string, string, error) {
	return ds.sh.Exec("docker stack deploy -c", file, stack)
}
