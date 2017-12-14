package cmd

type Gpm struct{}

func NewGpm() Gpm {
	return Gpm{}
}

func (Gpm) Install(installDir string, yaml string) (string, error) {
	commands := []string{"gpm install -F -c Containerfile"}
	if installDir != "" {
		commands = append(commands, "-d", installDir)
	}
	if yaml != "" {
		commands = append(commands, "-y", yaml)
	}
	return Sh().Exec(commands...)
}
