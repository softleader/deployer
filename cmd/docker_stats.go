package cmd

import (
	"github.com/softleader/deployer/models"
	"strings"
)

type DockerStats struct {
}

func NewDockerStats() *DockerStats {
	return &DockerStats{}
}

func (ds *DockerStats) NoStream(grep string) (s []models.DockerStatsNoStream, err error) {
	_, out, err := noStream(grep)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerStatsNoSteam(line))
		}
	}
	return
}

func noStream(grep string) (arg string, out string, err error) {
	commands := []string{"docker stats --no-stream --format '{{.Name}};{{.CPUPerc}};{{.MemUsage}};{{.MemPerc}};{{.NetIO}};{{.BlockIO}}'"}
	if grep != "" {
		commands = append(commands, "| grep", grep)
	}
	return Exec(&Options{}, commands...)
}
