package cmd

import (
	"github.com/softleader/deployer/datamodels"
	"strings"
)

type GenYaml struct {
	sh  Shell
	cmd string
}

func NewGenYaml(sh Shell, cmd string) *GenYaml {
	if cmd == "" {
		cmd = "gen-yaml"
	}
	return &GenYaml{sh: sh, cmd: cmd}
}

func (gy *GenYaml) Gen(opts *Options, output string, d *datamodels.Deploy, dirs ...string) (string, string, error) {

	commands := []string{gy.cmd, "-s swarm -o", output}
	if d.Silently {
		commands = append(commands, "-S")
	}
	if d.Net0 != "" {
		commands = append(commands, "--net0", d.Net0)
	}
	if d.Volume0 != "" {
		commands = append(commands, "--volume0", d.Volume0)
	}
	if d.Dev.Hostname != "" {
		commands = append(commands, d.Dev.String())
	}
	commands = append(commands, strings.Join(dirs, " "))

	return gy.sh.Exec(opts, commands...)
}

func (gy *GenYaml) Version() (string, string, error) {
	return gy.sh.Exec(&Options{}, gy.cmd, "--version")
}
