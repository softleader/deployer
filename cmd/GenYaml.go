package cmd

import (
	"github.com/softleader/deployer/datamodels"
	"strings"
)

type GenYaml struct {
	cmd string
}

func NewGenYaml(cmd string) GenYaml {
	if cmd == "" {
		cmd = "gen-yaml"
	}
	return GenYaml{cmd}
}

func (g GenYaml) Gen(output string, d datamodels.Deploy, dirs ...string) (string, error) {

	commands := []string{g.cmd, "-s swarm -o", output}
	if d.Silence {
		commands = append(commands, "-S")
	}
	if d.Net0 != "" {
		commands = append(commands, "--net0", d.Net0)
	}
	if d.Volume0 != "" {
		commands = append(commands, "--volume0", d.Volume0)
	}
	if d.Dev != "" {
		commands = append(commands, "--dev", d.Dev)
	}
	commands = append(commands, strings.Join(dirs, " "))

	return Sh().Exec(commands...)
}

func (g GenYaml) Version() (string, error) {
	return Sh().Exec(g.cmd, "-V")
}
