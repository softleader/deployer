package cmd

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

func (g *Gpm) Install(opts *Options, dir string, yaml string) (string, string, error) {
	commands := []string{g.cmd, "install -F -c Containerfile"}
	if dir != "" {
		commands = append(commands, "-d", dir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return g.sh.Exec(opts, commands...)
}

func (g *Gpm) Version() (string, string, error) {
	return g.sh.Exec(&Options{}, g.cmd, "--version")
}
