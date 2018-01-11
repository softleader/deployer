package cmd

type Gpm struct {
	sh  Shell
	cmd string
}

func NewGpm(sh Shell, cmd string) *Gpm {
	if cmd == "" {
		cmd = "gpm"
	}
	return &Gpm{sh: sh, cmd: cmd}
}

func (g *Gpm) Install(opts *Options, dir string, yaml string) (arg string, out string, err error) {
	commands := []string{g.cmd, "install -F -c Containerfile"}
	if dir != "" {
		commands = append(commands, "-d", dir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return g.sh.Exec(opts, commands...)
}

func (g *Gpm) Version() (arg string, out string, err error) {
	return g.sh.Exec(&Options{}, g.cmd, "--version")
}
