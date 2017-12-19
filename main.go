package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
	"os"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"log"
	"flag"
)

type args struct {
	wd      string
	addr    string
	port    int
	gpm     string
	genYaml string
}

func main() {
	args := newArgs()

	service := newService(args)
	checkDependencies(service)

	// https://github.com/kataras/iris
	serve(args, service)
}

func newArgs() args {
	a := args{}
	flag.StringVar(&a.wd, "wd", "", "Determine a working dictionary, default: $(pwd)/wd")
	flag.StringVar(&a.addr, "addr", "", " Determine application addr, default: empty")
	flag.IntVar(&a.port, "port", 5678, "Determine application port, default: 5678")
	flag.StringVar(&a.gpm, "cmd.gpm", "", "Command to execute softleader/git-package-manager, default: gpm")
	flag.StringVar(&a.genYaml, "cmd.gen-yaml", "", "Command to execute softleader/container-yaml-generator, default: gen-yaml")
	flag.Parse()
	return a
}

func newService(args args) services.DeployService {
	cmdWd := cmd.NewWd(args.wd)
	cmdSh := cmd.NewSh(cmdWd)
	return services.DeployService{
		DockerStack: cmd.NewDockerStack(cmdSh),
		Gpm:         cmd.NewGpm(cmdSh, args.gpm),
		GenYaml:     cmd.NewGenYaml(cmdSh, args.genYaml),
		Wd:          cmdWd,
	}
}

func checkDependencies(s services.DeployService) {
	fmt.Println("Checking dependencies...")
	cmd, out, err := s.Gpm.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
	cmd, out, err = s.GenYaml.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
}

func serve(args args, s services.DeployService) {
	app := iris.New()
	app.Controller("/", new(controller.DeployController), s)
	app.Run(
		iris.Addr(args.addr+":"+strconv.Itoa(args.port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}
