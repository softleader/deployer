package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/web/controller"
	"github.com/softleader/deployer/services"
	"os"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"log"
	"flag"
	"github.com/softleader/deployer/datamodels"
	"path"
)

type args struct {
	ws      string
	addr    string
	port    int
	gpm     string
	genYaml string
}

func main() {
	args := newArgs()

	service := newService(args)
	checkDependencies(*service)

	// https://github.com/kataras/iris
	serve(*args, *service)
}

func newArgs() *args {
	a := args{}
	flag.StringVar(&a.ws, "workspace", "", "Determine a workspace (default $(pwd)/workspace)")
	flag.StringVar(&a.addr, "addr", "", " Determine application address (default blank)")
	flag.IntVar(&a.port, "port", 5678, "Determine application port")
	flag.StringVar(&a.gpm, "cmd.gpm", "gpm", "Command to execute softleader/git-package-manager")
	flag.StringVar(&a.genYaml, "cmd.gen-yaml", "gen-yaml", "Command to execute softleader/container-yaml-generator")
	flag.Parse()
	return &a
}

func newService(args *args) *services.DeployService {
	ws := cmd.NewWorkspace(args.ws)
	sh := cmd.NewShell()
	return &services.DeployService{
		DockerStack: *cmd.NewDockerStack(*sh),
		Gpm:         *cmd.NewGpm(*sh, args.gpm),
		GenYaml:     *cmd.NewGenYaml(*sh, args.genYaml),
		Workspace:   *ws,
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

func serve(args args, s services.DeployService) {
	app := iris.New()

	tmpl := iris.HTML("web/views", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	app.Get("/deploy", func(ctx iris.Context) {
		ctx.ViewData("workspace", args.ws)
		ctx.ViewData("dft", datamodels.Deploy{
			Dev: datamodels.Dev{
				Hostname: "192.168.1.60",
				Port:     0,
				Ignore:   "elasticsearch,kibana,logstash,redis,eureka,softleader-config-server",
			},
			Yaml:    "github:softleader/softleader-package/",
			Volume0: "",
			Net0:    "",
			Group:   "",
		})
		ctx.View("deploy.html")
	})

	app.Controller("/", new(controller.StackController), s)
	app.Controller("/services", new(controller.ServiceController), s)

	app.Get("/download/{project:string}", func(ctx iris.Context) {
		pj := ctx.Params().Get("project")
		zip := s.Workspace.GetWd(false, pj).GetCompressPath()
		ctx.SendFile(zip, pj+"-"+path.Base(zip))
	})

	app.Run(
		iris.Addr(args.addr+":"+strconv.Itoa(args.port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}
