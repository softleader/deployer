package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/cache"
	"github.com/kataras/iris/context"
	"github.com/softleader/deployer/cmd"
)

func NewApplication(args *Args) *iris.Application {
	ws := NewWorkspace(args.Ws)
	commands := newCommands(args)

	// prepare routes
	stackRoutes := StackRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	serviceRoutes := ServiceRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	dashboardRoutes := DashboardRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	deployRoutes := DeployRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	practiceRoutes := PracticeRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	historyRoutes := HistoryRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}
	statsRoutes := StatsRoutes{Routes: Routes{Args: args, Workspace: ws, Commands: commands}}

	config := ws.Config
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
			ctx.Redirect(config.Index)
		})
	}

	dashboard := app.Party("/dashboard")
	{
		dashboardCache := cache.Handler(config.DashboardCache)
		dashboard.Get("/", dashboardRoutes.DashboardPage)
		dashboard.Get("/nodes", dashboardCache, dashboardRoutes.Nodes)
		dashboard.Get("/services", dashboardCache, dashboardRoutes.Services)
		dashboard.Get("/projects", dashboardCache, dashboardRoutes.Projects)
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

	stats := app.Party("/stats")
	{
		stats.Get("/", statsRoutes.GetStats)
	}

	return app
}

func newCommands(args *Args) *Commands {
	return &Commands{
		DockerStack:   *cmd.NewDockerStack(args.Registry.Login()),
		DockerService: *cmd.NewDockerService(),
		DockerNode:    *cmd.NewDockerNode(),
		Gpm:           *cmd.NewGpm(args.CmdGpm),
		GenYaml:       *cmd.NewGenYaml(args.CmdGenYaml),
	}
}
