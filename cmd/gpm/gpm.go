package gpm

import (
	"strings"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/cmd"
)

var Cmd string

func command() string {
	if Cmd == "" {
		return "gpm"
	}
	return Cmd
}

func Install(opts *cmd.Options, dir string, d *models.Deploy) (bool, error) {
	_, out, err := install(opts, dir, d.Yaml, d.Extend)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func install(opts *cmd.Options, dir string, yaml string, extend string) (arg string, out string, err error) {
	commands := []string{command(), "install -F -c Containerfile"}
	if dir != "" {
		commands = append(commands, "-d", dir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	if extend != "" {
		commands = append(commands, "-e", extend)
	}
	return cmd.Exec(opts, commands...)
}

func Version() (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, command(), "--version")
}
