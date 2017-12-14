package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
	"os"
)

// https://github.com/kataras/iris
func main() {

	if len(os.Args) > 1 {
		os.Chdir(os.Args[1])
	}

	app := iris.New()

	app.Controller("/", new(controller.DeploymentController), services.NewDeploymentService())

	app.Run(
		iris.Addr(":5678"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)

}
