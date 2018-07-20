package main

import (
	"github.com/kataras/iris"
	"strconv"
	"github.com/softleader/deployer/app"
)

var (
	Debug bool
)

func main() {
	args := app.NewArgs()

	// https://github.com/kataras/iris
	app.NewApplication(args, Debug).Run(
		iris.Addr(args.Addr+":"+strconv.Itoa(args.Port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}
