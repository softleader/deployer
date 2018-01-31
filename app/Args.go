package app

import "flag"

type Args struct {
	Ws      string
	Addr    string
	Port    int
	Gpm     string
	GenYaml string
	Debug   bool
}

func NewArgs() *Args {
	a := Args{}
	flag.StringVar(&a.Ws, "workspace", "", "Determine a workspace (default $(pwd)/workspace)")
	flag.StringVar(&a.Addr, "addr", "", " Determine application address (default blank)")
	flag.IntVar(&a.Port, "port", 5678, "Determine application port")
	flag.StringVar(&a.Gpm, "cmd.gpm", "gpm", "Command to execute softleader/git-package-manager")
	flag.StringVar(&a.GenYaml, "cmd.gen-yaml", "gen-yaml", "Command to execute softleader/container-yaml-generator")
	flag.BoolVar(&a.Debug, "debug", false, "Print logs to standard output.")
	flag.Parse()
	return &a
}
