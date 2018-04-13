package app

import (
	"flag"
	"fmt"
)

type Args struct {
	Ws       string
	Addr     string
	Port     int
	Gpm      string
	GenYaml  string
	Debug    bool
	Registry Registry
}

type Registry struct {
	Server   string
	Username string
	Password string
}

func NewArgs() *Args {
	a := Args{
		Registry: Registry{},
	}
	flag.StringVar(&a.Ws, "workspace", "", "Determine a workspace (default $(pwd)/workspace)")
	flag.StringVar(&a.Addr, "addr", "", " Determine application address (default blank)")
	flag.IntVar(&a.Port, "port", 5678, "Determine application port")
	flag.StringVar(&a.Registry.Server, "registry.server", "hub.softleader.com.tw", "SoftLeader Docker Registry")
	flag.StringVar(&a.Registry.Username, "registry.username", "client", "Username to SoftLeader Docker Registry")
	flag.StringVar(&a.Registry.Password, "registry.password", "poweredbysoftleader", "Password to SoftLeader Docker Registry")
	flag.StringVar(&a.Gpm, "cmd.gpm", "gpm", "Command to execute softleader/git-package-manager")
	flag.StringVar(&a.GenYaml, "cmd.gen-yaml", "gen-yaml", "Command to execute softleader/container-yaml-generator")
	flag.BoolVar(&a.Debug, "debug", false, "Print logs to standard output.")
	flag.Parse()
	return &a
}

func (r *Registry) Login() string {
	return fmt.Sprintf("docker login %s -u %s -p %s", r.Server, r.Username, r.Password)
}
