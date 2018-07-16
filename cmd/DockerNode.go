package cmd

type DockerNode struct {
}

func NewDockerNode() *DockerNode {
	return &DockerNode{}
}

func (dn *DockerNode) Ls() (arg string, out string, err error) {
	return Exec(&Options{}, "docker node ls", "--format '{{.Availability}}'")
}
