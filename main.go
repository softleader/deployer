package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
	"os"
	"flag"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"log"
)

// https://github.com/kataras/iris
func main() {

	wd := flag.String("wd", "", "Determine a working dictionary, default: $(pwd)/wd")
	addr := flag.String("addr", "", " Determine application addr, default: empty")
	port := flag.Int("port", 5678, "Determine application port, default: 5678")
	gpm := flag.String("cmd.gpm", "", "Command to execute softleader/git-package-manager, default: gpm")
	genYaml :=
		flag.String("cmd.gen-yaml", "", "Command to execute softleader/container-yaml-generator, default: gen-yaml")

	flag.Parse()

	fmt.Printf("Setting up working directory to '%v'\n", *wd)
	cmdWd := cmd.NewWd(*wd)

	cmdSh := cmd.NewSh(cmdWd)

	s := services.DeployService{
		DockerStack: cmd.NewDockerStack(cmdSh),
		Gpm:         cmd.NewGpm(cmdSh, *gpm),
		GenYaml:     cmd.NewGenYaml(cmdSh, *genYaml),
		Wd:          cmdWd,
	}

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

	serve(*addr+":"+strconv.Itoa(*port), s)
}

func serve(addr string, s services.DeployService) {
	app := iris.New()

	app.Controller("/", new(controller.DeployController), s)

	app.Run(
		iris.Addr(addr),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}
