package cmd

import (
	"github.com/softleader/deployer/datamodels"
	"strings"
	"fmt"
)

type GenYaml struct{}

func NewGenYaml() GenYaml {
	return GenYaml{}
}

func (GenYaml) Gen(output string, d datamodels.Deployment, dirs ...string) (string, error) {

	commands := []string{"gen-yaml -s swarm -o", output}
	if d.Net0 != "" {
		commands = append(commands, "--net0", d.Net0)
	}
	if d.Volume0 != "" {
		commands = append(commands, "--volume0", d.Volume0)
	}
	commands = append(commands, "--dev", fmt.Sprintf("192.168.1.60/%v", d.PublishPort))
	commands = append(commands, strings.Join(dirs, " "))

	return Sh().Exec(commands...)
}
