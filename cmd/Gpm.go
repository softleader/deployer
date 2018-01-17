package cmd

import (
	"github.com/softleader/deployer/pipe"
	"strings"
	"github.com/softleader/deployer/models"
	"os"
	"fmt"
	"github.com/labstack/gommon/log"
)

type Gpm struct {
	cmd string
}

func NewGpm(cmd string) *Gpm {
	if cmd == "" {
		cmd = "gpm"
	}
	gpm := Gpm{cmd: cmd}

	cmd, out, err := gpm.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
	return &gpm
}

func (g *Gpm) Install(opts *Options, dir string, d *models.Deploy) (bool, error) {
	_, out, err := install(g.cmd, opts, dir, d.Yaml)
	if err != nil {
		return false, err
	}
	(*opts.Ctx).StreamWriter(pipe.Print(out))
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func install(cmd string, opts *Options, dir string, yaml string) (arg string, out string, err error) {
	commands := []string{cmd, "install -F -c Containerfile"}
	if dir != "" {
		commands = append(commands, "-d", dir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return Exec(opts, commands...)
}

func (g *Gpm) Version() (arg string, out string, err error) {
	return Exec(&Options{}, g.cmd, "--version")
}
