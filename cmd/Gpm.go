package cmd

type Gpm struct {
	sh  Sh
	cmd string
}

func NewGpm(sh Sh, cmd string) Gpm {
	if cmd == "" {
		cmd = "gpm"
	}
	return Gpm{sh: sh, cmd: cmd}
}

func (g Gpm) Install(installDir string, yaml string) (string, string, error) {
	commands := []string{g.cmd, "install -F -c Containerfile"}
	if installDir != "" {
		commands = append(commands, "-d", installDir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return g.sh.Exec(commands...)
}

func (g Gpm) Version() (string, string, error) {
	return g.sh.Exec(g.cmd, "--version")
}
