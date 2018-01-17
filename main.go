package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/services"
	"os"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/cmd"
	"log"
	"flag"
	"github.com/softleader/deployer/models"
	"path"
	"strings"
	"time"
	"github.com/softleader/deployer/pipe"
	"encoding/json"
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

	ds := newDeployService(args)
	checkDependencies(*ds)

	ps := newPracticeService(args)

	// https://github.com/kataras/iris
	app := newApp(*args, *ds, *ps)

	app.Run(
		iris.Addr(args.addr+":"+strconv.Itoa(args.port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
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

func newPracticeService(args *args) *services.PracticeService {
	ws := cmd.NewWorkspace(args.ws)
	return &services.PracticeService{
		Workspace: *ws,
	}
}

func newDeployService(args *args) *services.DeployService {
	ws := cmd.NewWorkspace(args.ws)
	sh := cmd.NewShell()
	return &services.DeployService{
		DockerStack:   *cmd.NewDockerStack(*sh),
		DockerService: *cmd.NewDockerService(*sh),
		Gpm:           *cmd.NewGpm(*sh, args.gpm),
		GenYaml:       *cmd.NewGenYaml(*sh, args.genYaml),
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

func newApp(args args, ds services.DeployService, ps services.PracticeService) *iris.Application {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	deployRoutes := app.Party("/deploy")
	{
		deployRoutes.Get("/", func(ctx iris.Context) {
			ctx.ViewData("workspace", args.ws)
			ctx.ViewData("dft", models.Deploy{
				Dev: models.Dev{
					IpAddress: "192.168.1.60",
					Port:      0,
					Ignore:    "elasticsearch,kibana,logstash,redis,eureka,softleader-config-server",
				},
				Yaml:    "github:softleader/softleader-package/",
				Volume0: "",
				Net0:    "",
				Group:   "",
			})
			ctx.View("deploy.html")
		})

		deployRoutes.Get("/download/{project:string}", func(ctx iris.Context) {
			pj := ctx.Params().Get("project")
			zip := ds.Workspace.GetWd(false, pj).GetCompressPath()
			ctx.SendFile(zip, pj+"-"+path.Base(zip))
		})
	}

	stacksRoutes := app.Party("/")
	{
		stacksRoutes.Get("/", func(ctx iris.Context) {
			out, err := ds.GetAll()
			if err != nil {
				out = append(out, []string{err.Error()})
			}

			stacks := make(map[string][][]string)
			for _, line := range out {
				splited := strings.Split(line[0], "-")
				key := splited[0]
				if len(splited) > 1 {
					if publishedPort(splited[1]) { // 有 publish port 可視為有開啟 dev 模式
						key = strings.Join(splited[:2], "-")
					}
				}
				_, out, _ := ds.DockerService.GetCreatedTimeOfFirstServiceInStack(line[0])
				out = strings.TrimSuffix(out, "\n")
				t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", out)
				line = append(line, uptime(t))
				stacks[key] = append(stacks[key], line)
			}
			ctx.ViewData("stacks", stacks)
			ctx.View("stack.html")
		})

		stacksRoutes.Post("/", func(ctx iris.Context) {
			d := &models.Deploy{}
			ctx.ReadJSON(d)
			start := time.Now()
			indent, _ := json.MarshalIndent(d, "", " ")

			ctx.StreamWriter(pipe.Printf("Received deploy request: %v", string(indent)))
			err := ds.Deploy(&ctx, *d)
			if err != nil {
				ctx.Application().Logger().Warn(err.Error())
				ctx.WriteString(err.Error())
			}
			ctx.StreamWriter(pipe.Printf("Resolving in %v, done.", time.Since(start)))
		})

		stacksRoutes.Get("/rm/{stack:string}", func(ctx iris.Context) {
			stack := ctx.Params().Get("stack")
			_, err := ds.DeleteStack(stack)
			if err != nil {
				ctx.Application().Logger().Warn(err.Error())
				ctx.WriteString(err.Error())
			}
			ctx.Redirect("/")
		})
	}

	servicesRoutes := app.Party("/services")
	{
		servicesRoutes.Get("/{stack:string}", func(ctx iris.Context) {
			stack := ctx.Params().Get("stack")
			out, err := ds.GetServices(stack)
			if err != nil {
				out = append(out, []string{err.Error()})
			}
			ctx.ViewData("out", out)
			ctx.ViewData("stack", stack)
			ctx.View("service.html")
		})

		servicesRoutes.Get("/ps/{serviceId:string}", func(ctx iris.Context) {
			serviceId := ctx.Params().Get("serviceId")
			out, err := ds.Ps(serviceId)
			if err != nil {
				ctx.Application().Logger().Warn(err.Error())
				ctx.WriteString(err.Error())
			}
			ctx.ViewData("out", out)
			ctx.View("ps.html")
		})

		servicesRoutes.Get("/rm/{stack:string}/{service:string}", func(ctx iris.Context) {
			stack := ctx.Params().Get("stack")
			service := ctx.Params().Get("service")
			_, err := ds.DeleteService(service)
			if err != nil {
				ctx.Application().Logger().Warn(err.Error())
				ctx.WriteString(err.Error())
			}
			ctx.Redirect("/services/" + stack)
		})

	}

	practicesRoutes := app.Party("/best-practices")
	{
		practicesRoutes.Get("/", func(ctx iris.Context) {
			out, err := ps.Get()
			ctx.ViewData("err", err)
			ctx.ViewData("out", out)
			ctx.View("best-practices.html")
		})

		practicesRoutes.Get("/mde", func(ctx iris.Context) {
			out, err := ps.Get()
			ctx.ViewData("err", err)
			ctx.ViewData("out", out)
			ctx.View("best-practices-mde.html")
		})

		practicesRoutes.Post("/mde", func(ctx iris.Context) {
			c := ctx.PostValue("content");
			err := ps.Save(c)
			if err != nil {
				ctx.Application().Logger().Warn(err.Error())
				ctx.WriteString(err.Error())
			}
			ctx.Redirect("/best-practices")
		})
	}

	return app
}

func uptime(t time.Time) string {
	d := time.Since(t)
	return fmt.Sprintf("up %.2f hours", d.Hours())
}

func publishedPort(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}
