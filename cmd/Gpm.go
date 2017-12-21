package cmd

import "github.com/kataras/iris"

type Gpm struct {
	sh  Sh
	cmd string
}

func NewGpm(sh Sh, cmd string) *Gpm {
	if cmd == "" {
		cmd = "gpm"
	}
	return &Gpm{sh: sh, cmd: cmd}
}

func (g *Gpm) Install(ctx *iris.Context, dir string, yaml string) (string, string, error) {
	commands := []string{g.cmd, "install -F -c Containerfile"}
	if dir != "" {
		commands = append(commands, "-d", dir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return g.sh.Exec(ctx, commands...)
}

func (g *Gpm) Version() (string, string, error) {
	return g.sh.Exec(nil, g.cmd, "--version")
}
