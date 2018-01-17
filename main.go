package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/services"
	"os"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"log"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/routes"
)

func main() {
	args := models.NewArgs()

	ds := newDeployService(args)
	checkDependencies(*ds)

	ps := newPracticeService(args)

	// https://github.com/kataras/iris
	app := newApp(*args, *ds, *ps)

	app.Run(
		iris.Addr(args.Addr+":"+strconv.Itoa(args.Port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}

func newPracticeService(args *models.Args) *services.PracticeService {
	ws := cmd.NewWorkspace(args.Ws)
	return &services.PracticeService{
		Workspace: *ws,
	}
}

func newDeployService(args *models.Args) *services.DeployService {
	ws := cmd.NewWorkspace(args.Ws)
	sh := cmd.NewShell()
	return &services.DeployService{
		DockerStack:   *cmd.NewDockerStack(*sh),
		DockerService: *cmd.NewDockerService(*sh),
		Gpm:           *cmd.NewGpm(*sh, args.Gpm),
		GenYaml:       *cmd.NewGenYaml(*sh, args.GenYaml),
		Workspace:     *ws,
	}
}

func checkDependencies(s services.DeployService) {
	fmt.Println("Checking dependencies...")
	cmd, out, err := s.Gpm.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
	cmd, out, err = s.GenYaml.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
}

func newApp(args models.Args, ds services.DeployService, ps services.PracticeService) *iris.Application {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	deployRoutes := routes.NewDeployRoutes(args, ds)
	deploy := app.Party("/deploy")
	{
		deploy.Get("/", deployRoutes.DeployPage)
		deploy.Get("/download/{project:string}", deployRoutes.DownloadYAML)
	}

	stackRoutes := routes.NewStackRoutes(args, ds)
	stacks := app.Party("/")
	{
		stacks.Get("/", stackRoutes.ListStack)
		stacks.Post("/", stackRoutes.DeployStack)
		stacks.Get("/rm/{stack:string}", stackRoutes.RemoveStack)
	}

	serviceRoutes := routes.NewServiceRoutes(args, ds)
	services := app.Party("/services")
	{
		services.Get("/{stack:string}", serviceRoutes.ListService)
		services.Get("/ps/{serviceId:string}", serviceRoutes.PsService)
		services.Get("/rm/{stack:string}/{service:string}", serviceRoutes.RemoveService)
	}

	practiceRoutes := routes.NewPracticeRoutes(args, ps)
	practices := app.Party("/best-practices")
	{
		practices.Get("/", practiceRoutes.BestPractices)
		practices.Get("/mde", practiceRoutes.MarkdownEditor)
		practices.Post("/mde", practiceRoutes.SaveMarkdown)
	}

	return app
}
