package docker

import (
	"github.com/softleader/deployer/models"
	"strings"
	"fmt"
	"sync"
	"github.com/softleader/deployer/cmd"
)

func StatsNoStream(grep string) (s []models.DockerStatsNoStream, err error) {
	nodes, err := readyNodes()
	if err != nil {
		return
	}
	out, err := parallelOverNodes(grep, nodes, sshNoStream)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, models.NewDockerStatsNoSteam(line))
		}
	}
	return
}

// this function make test possible
func parallelOverNodes(grep string, nodes []string, consume func(grep string, host string) string) (out string, err error) {
	c := make(chan string, len(nodes))
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			c <- consume(grep, node)
		}(node)
	}
	wg.Wait()
	close(c)
	for o := range c {
		out += o
	}
	return
}

func readyNodes() (nodes []string, err error) {
	n, err := NodeLs()
	for _, node := range n {
		if node.Status == "Ready" {
			nodes = append(nodes, node.Hostname)
		}
	}
	return
}

func sshNoStream(grep string, host string) (out string) {
	commands := []string{
		fmt.Sprintf(`ssh %s "docker stats --no-stream --format '{{.Name}};{{.CPUPerc}};{{.MemUsage}};{{.MemPerc}};{{.NetIO}};{{.BlockIO}}'`, host)}
	if grep != "" {
		commands = append(commands, "| grep", grep)
	}
	commands = append(commands, `"`)
	_, out, e := cmd.Exec(&cmd.Options{}, commands...)
	if e != nil {
		fmt.Println(e)
	}
	return
}
