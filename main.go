package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
)

// https://github.com/kataras/iris
func main() {
	app := iris.New()

	app.Controller("/", new(controller.DeploymentController), new(services.DeploymentService))

	app.Run(
		iris.Addr(":5678"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)

}
