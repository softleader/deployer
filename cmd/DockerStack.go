package cmd

import "github.com/kataras/iris"

type DockerStack struct {
	sh Sh
}

func NewDockerStack(sh Sh) *DockerStack {
	return &DockerStack{sh: sh}
}

func (ds *DockerStack) Ls() (string, string, error) {
	return ds.sh.Exec("docker stack ls")
}

func (ds *DockerStack) Services(stack string) (string, string, error) {
	return ds.sh.Exec("docker stack services", stack)
}

func (ds *DockerStack) Ps(id string) (string, string, error) {
	return ds.sh.Exec("docker service ps", id, "--no-trunc")
}

func (ds *DockerStack) Rm(stack string) (string, string, error) {
	return ds.sh.Exec("docker stack rm", stack)
}

func (ds *DockerStack) RmLike(stack string) (string, string, error) {
	return ds.sh.Exec("docker stack rm $(docker stack ls --format {{.Name}} | grep", stack, ")")
}

func (ds *DockerStack) Deploy(ctx *iris.Context, stack string, file string) {
	ds.sh.ExecPipe(ctx, "docker stack deploy -c", file, stack)
}
