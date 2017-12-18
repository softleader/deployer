package cmd

type Gpm struct {
	cmd string
}

func NewGpm(cmd string) Gpm {
	if cmd == "" {
		cmd = "gpm"
	}
	return Gpm{cmd}
}

func (g Gpm) Install(installDir string, yaml string) (string, string, error) {
	commands := []string{g.cmd, "install -F -c Containerfile"}
	if installDir != "" {
		commands = append(commands, "-d", installDir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return Sh().Exec(commands...)
}


func (g Gpm) Version() (string, string, error) {
	return Sh().Exec(g.cmd, "--version")
}
