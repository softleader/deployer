package docker

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/models"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"strings"
	"sync"
)

var (
	cache, _ = homedir.Expand("~/.config/dockerctl")
)

func StatsNoStream(grep, cache string) (s []models.DockerStatsNoStream, err error) {
	nodes, err := readyNodes(cache)
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
	n, err := dockerd.LoadNodes(logrus.StandardLogger(), cache, false)
	if err == nil {
		for _, stats := range s {
			stats.Addr = n[stats.Name].Addr
		}
	}
	return
}

// this function make test possible
func parallelOverNodes(grep string, nodes []dockerd.Node, consume func(grep string, node dockerd.Node) string) (out string, err error) {
	c := make(chan string, len(nodes))
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(node dockerd.Node) {
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

func readyNodes(cache string) ([]dockerd.Node, error) {
	nodes, err := dockerd.LoadNodes(logrus.StandardLogger(), cache, false)
	if err != nil {
		return nil, err
	}
	var ns []dockerd.Node
	for _, n := range nodes {
		ns = append(ns, n)
	}
	return ns, nil
}

func sshNoStream(grep string, node dockerd.Node) (out string) {
	commands := []string{
		fmt.Sprintf(`ssh %s "docker stats --no-stream --format '{{.Name}};{{.CPUPerc}};{{.MemUsage}};{{.MemPerc}};{{.NetIO}};{{.BlockIO}}'`, node.Addr)}
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
