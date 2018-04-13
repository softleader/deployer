package main

import (
	"github.com/kataras/iris"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/routes"
	"github.com/softleader/deployer/app"
	"github.com/kataras/iris/context"
)

var (
	Debug bool
)

func main() {
	args := app.NewArgs()
	ws := app.NewWorkspace(args.Ws)

	deployRoutes := newDeployRoutes(args, ws)
	serviceRoutes := newServiceRoutes(args)
	stackRoutes := newStackRoutes(args, ws)
	practiceRoutes := newPracticeRoutes(ws)
	historyRoutes := newHistoryRoutes(ws)

	// https://github.com/kataras/iris
	app := newApp(deployRoutes, stackRoutes, serviceRoutes, practiceRoutes, historyRoutes)

	app.Run(
		iris.Addr(args.Addr+":"+strconv.Itoa(args.Port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}

func newDeployRoutes(args *app.Args, ws *app.Workspace) *routes.DeployRoutes {
	return &routes.DeployRoutes{
		Args:      *args,
		Workspace: *ws,
	}
}

func newServiceRoutes(args *app.Args) *routes.ServiceRoutes {
	return &routes.ServiceRoutes{
		DockerStack:   *cmd.NewDockerStack(args.Registry),
		DockerService: *cmd.NewDockerService(),
	}
}

func newStackRoutes(args *app.Args, ws *app.Workspace) *routes.StackRoutes {
	return &routes.StackRoutes{
		Workspace:     *ws,
		DockerStack:   *cmd.NewDockerStack(args.Registry),
		DockerService: *cmd.NewDockerService(),
		Gpm:           *cmd.NewGpm(args.Gpm),
		GenYaml:       *cmd.NewGenYaml(args.GenYaml),
		Debug:         args.Debug,
	}
}

func newHistoryRoutes(ws *app.Workspace) *routes.HistoryRoutes {
	return &routes.HistoryRoutes{
		Workspace: *ws,
	}
}

func newPracticeRoutes(ws *app.Workspace) *routes.PracticeRoutes {
	return &routes.PracticeRoutes{
		Workspace: *ws,
	}
}

func newApp(deployRoutes *routes.DeployRoutes,
	stackRoutes *routes.StackRoutes,
	serviceRoutes *routes.ServiceRoutes,
	practiceRoutes *routes.PracticeRoutes,
	historyRoutes *routes.HistoryRoutes) *iris.Application {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	api := app.Party("/api")
	{
		api.Post("/stacks", stackRoutes.DeployStack)
		api.Delete("/stacks/{stack:string}", stackRoutes.RemoveStack)
		api.Delete("/services/{service:string}", serviceRoutes.RemoveService)
	}

	root := app.Party("/")
	{
		root.Get("/", deployRoutes.DeployPage)
		root.Get("/{history:string}", deployRoutes.DeployPage)
	}

	// deprecate, in place of root route
	deploy := app.Party("/deploy")
	{
		deploy.Get("/", func(ctx context.Context) {
			ctx.Redirect("/")
		})
		deploy.Get("/{history:string}", func(ctx context.Context) {
			h := ctx.Params().Get("history")
			ctx.Redirect("/" + h)
		})
	}

	yamls := app.Party("/yamls")
	{
		yamls.Get("/{project:string}", deployRoutes.DownloadYAML)
		yamls.Post("/", stackRoutes.GenerateYAML)
	}

	stacks := app.Party("/stacks")
	{
		stacks.Get("/", stackRoutes.ListStack)
		stacks.Post("/", stackRoutes.DeployStack)
		stacks.Get("/rm/{stack:string}", func(ctx context.Context) {
			stackRoutes.RemoveStack(ctx)
			ctx.Redirect("/stacks")
		})
	}

	services := app.Party("/services")
	{
		services.Get("/{stack:string}", serviceRoutes.ListService)
		services.Get("/ps/{serviceId:string}", serviceRoutes.PsService)
		services.Get("/rm/{stack:string}/{service:string}", func(ctx context.Context) {
			serviceRoutes.RemoveService(ctx)
			stack := ctx.Params().Get("stack")
			ctx.Redirect("/services/" + stack)
		})
	}

	practices := app.Party("/best-practices")
	{
		practices.Get("/", practiceRoutes.BestPractices)
		practices.Get("/mde", practiceRoutes.MarkdownEditor)
		practices.Post("/mde", practiceRoutes.SaveMarkdown)
	}

	history := app.Party("/history")
	{
		history.Get("/", historyRoutes.GetHistory)
		history.Get("/rm/{idx:int}", func(ctx context.Context) {
			historyRoutes.RemoveHistory(ctx)
			ctx.Redirect("/history")
		})
	}

	return app
}
