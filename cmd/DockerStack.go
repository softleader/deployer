package cmd

type DockerStack struct{}

func NewDockerStack() DockerStack {
	return DockerStack{}
}

func (DockerStack) Ls() (string, string, error) {
	return Sh().Exec("docker stack ls")
}

func (DockerStack) Services(stack string) (string, string, error) {
	return Sh().Exec("docker stack services", stack)
}

func (DockerStack) Rm(stack string) (string, string, error) {
	return Sh().Exec("docker stack rm", stack)
}

func (DockerStack) Deploy(stack string, file string) (string, string, error) {
	return Sh().Exec("docker stack deploy -c", file, stack)
}
