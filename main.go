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
)

// https://github.com/kataras/iris
func main() {

	wd := flag.String("wd", "", "Determine a working dictionary, default: $(pwd)")
	addr := flag.String("addr", "", " Determine application addr, default: empty")
	port := flag.Int("port", 5678, "Determine application port, default: 5678")
	gpm := flag.String("cmd.gpm", "", "Command to execute softleader/git-package-manager, default: gpm")
	genYaml :=
		flag.String("cmd.gen-yaml", "", "Command to execute softleader/container-yaml-generator, default: gen-yaml")

	flag.Parse()

	if *wd != "" {
		os.Chdir(*wd)
		pwd, _ := os.Getwd()
		fmt.Printf("Changed working directory to [%v]\n", pwd)
	}

	app := iris.New()

	service := services.DeployService{
		DockerStack: cmd.NewDockerStack(),
		Gpm:         cmd.NewGpm(*gpm),
		GenYaml:     cmd.NewGenYaml(*genYaml),
		Wd:          cmd.NewWd(),
	}

	app.Controller("/", new(controller.DeployController), service)

	app.Run(
		iris.Addr(*addr+":"+strconv.Itoa(*port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)

}
