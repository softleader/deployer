package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/cache"
	"github.com/kataras/iris/context"
	"github.com/mitchellh/go-homedir"
	"github.com/softleader/deployer/cmd/genYaml"
	"github.com/softleader/deployer/cmd/gpm"
)

var Ws *Workspace
var Args *Arguments
var Debug bool

func NewApplication(args *Arguments, debug bool) *iris.Application {
	Args = args
	Ws = NewWorkspace(args.Ws)
	Debug = debug

	injectGpmCommand()
	injectGenYamlCommand()

	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)
	app.StaticWeb("/", "./static")

	app.UseGlobal(func(ctx iris.Context) {
		ctx.ViewData("navbar", Ws.Config.Navbar)
		ctx.Next() // execute the next handler, in this case the main one.
	})

	api := app.Party("/api")
	{
		api.Post("/stacks", DeployStack)
		api.Delete("/stacks/{stack:string}", RemoveStack)
		api.Delete("/services/{service:string}", RemoveService)
		api.Put("/services/{serviceId:string}", UpdateService)
	}

	root := app.Party("/")
	{
		root.Get("/", func(ctx context.Context) {
			ctx.Redirect(Ws.Config.Index)
		})
	}

	dashboard := app.Party("/dashboard")
	{
		dashboardCache := cache.Handler(Ws.Config.DashboardCache)
		dashboard.Get("/", DashboardPage)
		dashboard.Get("/nodes", dashboardCache, DashboardNodes)
		dashboard.Get("/services", dashboardCache, DashboardServices)
		dashboard.Get("/projects", dashboardCache, DashboardProjects)
	}

	deploy := app.Party("/deploy")
	{
		deploy.Get("/", DeployPage)
		deploy.Get("/{history:string}", DeployPage)
	}

	yamls := app.Party("/yamls")
	{
		yamls.Get("/{project:string}", DownloadYAML)
		yamls.Post("/", GenerateYAML)
	}

	stacks := app.Party("/stacks")
	{
		stacks.Get("/", StackPage)
		stacks.Get("/ls", ListStack)
		stacks.Post("/", DeployStack)
		stacks.Get("/rm/{stack:string}", func(ctx context.Context) {
			RemoveStack(ctx)
			ctx.Redirect("/stacks")
		})
	}

	services := app.Party("/services")
	{
		services.Get("/{stack:string}", ListService)
		services.Get("/filter", FilterService)
		services.Get("/ps/{serviceId:string}", PsService)
		services.Get("/inspect/{serviceId:string}", InspectService)
		services.Get("/update/{serviceId:string}", UpdateService)
		services.Get("/logs/{serviceId:string}/{tail:int}", LogsService)
		services.Get("/rm/{stack:string}/{service:string}", func(ctx context.Context) {
			RemoveService(ctx)
			stack := ctx.Params().Get("stack")
			ctx.Redirect("/services/" + stack)
		})
	}

	practices := app.Party("/best-practices")
	{
		practices.Get("/", BestPractices)
		practices.Get("/mde", MarkdownEditor)
		practices.Post("/mde", SaveMarkdown)
	}

	history := app.Party("/histories")
	{
		history.Get("/", GetHistory)
		history.Get("/rm/{idx:int}", func(ctx context.Context) {
			RemoveHistory(ctx)
			ctx.Redirect("/histories")
		})
	}

	cache, _ := homedir.Expand(args.NodeCache)
	fmt.Printf("Setting up node stats cache location to '%v'\n", cache)
	stats := app.Party("/stats")
	{
		stats.Get("/", func(ctx context.Context) {
			GetStats(ctx, cache)
		})
	}

	return app
}

func injectGpmCommand() {
	gpm.Cmd = Args.Gpm
	cmd, out, err := gpm.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
}

func injectGenYamlCommand() {
	genYaml.Cmd = Args.GenYaml
	cmd, out, err := genYaml.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
}
