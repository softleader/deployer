package cmd

type DockerStats struct {
}

func NewDockerStats() *DockerStats {
	return &DockerStats{}
}

func (ds *DockerStats) NoStream(grep string) (arg string, out string, err error) {
	commands := []string{"docker stats --no-stream --format '{{.Name}};{{.CPUPerc}};{{.MemUsage}};{{.MemPerc}}'"}
	if grep != "" {
		commands = append(commands, "| grep", grep)
	}
	return Exec(&Options{}, commands...)
}
