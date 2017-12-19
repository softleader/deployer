package cmd

import (
	"github.com/softleader/deployer/datamodels"
	"strings"
)

type GenYaml struct {
	sh  Sh
	cmd string
}

func NewGenYaml(sh Sh, cmd string) *GenYaml {
	if cmd == "" {
		cmd = "gen-yaml"
	}
	return &GenYaml{sh: sh, cmd: cmd}
}

func (gy *GenYaml) Gen(output string, d *datamodels.Deploy, dirs ...string) (string, string, error) {

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
	if d.Dev.Addr != "" {
		commands = append(commands, "--dev", d.Dev.String())
	}
	commands = append(commands, strings.Join(dirs, " "))

	return gy.sh.Exec(commands...)
}

func (gy *GenYaml) Version() (string, string, error) {
	return gy.sh.Exec(gy.cmd, "--version")
}
