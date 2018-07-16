package main

import (
	"github.com/kataras/iris"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/routes"
	"github.com/softleader/deployer/app"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/cache"
	"time"
)

var (
	Debug bool
)

func main() {
	args := app.NewArgs()
	ws := app.NewWorkspace(args.Ws)

	deployRoutes := newDeployRoutes(args, ws)
	serviceRoutes := newServiceRoutes(args, ws)
	stackRoutes := newStackRoutes(args, ws)
	practiceRoutes := newPracticeRoutes(ws)
	historyRoutes := newHistoryRoutes(ws)
	dashboardRoutes := newDashboardRoutes(args, ws)

	// https://github.com/kataras/iris
	app := newApp(deployRoutes, stackRoutes, serviceRoutes, practiceRoutes, historyRoutes, dashboardRoutes)

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

func newServiceRoutes(args *app.Args, ws *app.Workspace) *routes.ServiceRoutes {
	return &routes.ServiceRoutes{
		Workspace:     *ws,
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

func newDashboardRoutes(args *app.Args, ws *app.Workspace) *routes.DashboardRoutes {
	return &routes.DashboardRoutes{
		Workspace:     *ws,
		DockerNode:    *cmd.NewDockerNode(),
		DockerService: *cmd.NewDockerService(),
		DockerStack:   *cmd.NewDockerStack(args.Registry),
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
	historyRoutes *routes.HistoryRoutes,
	dashboardRoutes *routes.DashboardRoutes) *iris.Application {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	app.StaticWeb("/", "./static")

	api := app.Party("/api")
	{
		api.Post("/stacks", stackRoutes.DeployStack)
		api.Delete("/stacks/{stack:string}", stackRoutes.RemoveStack)
		api.Delete("/services/{service:string}", serviceRoutes.RemoveService)
	}

	root := app.Party("/")
	{
		root.Get("/", func(ctx context.Context) {
			ctx.Redirect(deployRoutes.Config.Index)
		})
	}

	dashboard := app.Party("/dashboard")
	{
		dashboard.Get("/", dashboardRoutes.DashboardPage)
		dashboard.Get("/nodes", dashboardRoutes.Nodes)
		dashboard.Get("/services", dashboardRoutes.Services)
		dashboard.Get("/projects", cache.Handler(3*time.Second), dashboardRoutes.Projects)
	}

	deploy := app.Party("/deploy")
	{
		deploy.Get("/", deployRoutes.DeployPage)
		deploy.Get("/{history:string}", deployRoutes.DeployPage)
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
		services.Get("/inspect/{serviceId:string}", serviceRoutes.InspectService)
		services.Get("/update/{serviceId:string}", serviceRoutes.UpdateService)
		services.Get("/logs/{serviceId:string}/{tail:int}", serviceRoutes.LogsService)
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

	history := app.Party("/histories")
	{
		history.Get("/", historyRoutes.GetHistory)
		history.Get("/rm/{idx:int}", func(ctx context.Context) {
			historyRoutes.RemoveHistory(ctx)
			ctx.Redirect("/histories")
		})
	}

	return app
}
