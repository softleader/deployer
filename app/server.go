package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/cache"
	"github.com/kataras/iris/context"
	"github.com/softleader/deployer/cmd"
)

func NewApplication(args *Args) *iris.Application {
	ws := NewWorkspace(args.Ws)
	r := newRoute(args, ws)

	// prepare routes
	stackRoutes := StackRoutes{Route: r}
	serviceRoutes := ServiceRoutes{Route: r}
	dashboardRoutes := DashboardRoutes{Route: r}
	deployRoutes := DeployRoutes{Route: r}
	practiceRoutes := PracticeRoutes{Route: r}
	historyRoutes := HistoryRoutes{Route: r}
	statsRoutes := StatsRoutes{Route: r}

	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)
	app.StaticWeb("/", "./static")

	app.UseGlobal(func(ctx iris.Context) {
		ctx.ViewData("navbar", ws.Config.Navbar)
		ctx.Next() // execute the next handler, in this case the main one.
	})

	api := app.Party("/api")
	{
		api.Post("/stacks", stackRoutes.DeployStack)
		api.Delete("/stacks/{stack:string}", stackRoutes.RemoveStack)
		api.Delete("/services/{service:string}", serviceRoutes.RemoveService)
	}

	root := app.Party("/")
	{
		root.Get("/", func(ctx context.Context) {
			ctx.Redirect(ws.Config.Index)
		})
	}

	dashboard := app.Party("/dashboard")
	{
		dashboardCache := cache.Handler(ws.Config.DashboardCache)
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

func newRoute(args *Args, ws *Workspace) *Route {
	return &Route{
		Args:      args,
		Workspace: ws,
		Gpm:       *cmd.NewGpm(args.CmdGpm),
		GenYaml:   *cmd.NewGenYaml(args.CmdGenYaml),
	}
}
