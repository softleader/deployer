package cmd

import (
	"github.com/softleader/deployer/datamodels"
	"strings"
)

type GenYaml struct{}

func NewGenYaml() GenYaml {
	return GenYaml{}
}

func (GenYaml) Gen(output string, d datamodels.Deploy, dirs ...string) (string, error) {

	commands := []string{"gen-yaml -s swarm -o", output}
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
