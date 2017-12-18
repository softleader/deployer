package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
	"os"
	"flag"
	"fmt"
	"strconv"
)

// https://github.com/kataras/iris
func main() {

	wd := flag.String("wd", "", "determine a working dictionary")
	addr := flag.String("addr", "", "addr")
	port := flag.Int("port", 5678, "port")

	flag.Parse()

	if *wd != "" {
		os.Chdir(*wd)
		pwd, _ := os.Getwd()
		fmt.Printf("Changed working directory to [%v]\n", pwd)
	}

	app := iris.New()

	app.Controller("/", new(controller.DeployController), services.NewDeployService())

	app.Run(
		iris.Addr(*addr+":"+strconv.Itoa(*port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)

}
