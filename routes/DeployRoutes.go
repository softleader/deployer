package routes

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"path"
	"github.com/softleader/deployer/app"
)

type DeployRoutes struct {
	app.Args
	app.Workspace
}

func (r *DeployRoutes) DeployPage(ctx iris.Context) {
	ctx.ViewData("workspace", r.Args.Ws)
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
}

func (r *DeployRoutes) DownloadYAML(ctx iris.Context) {
	pj := ctx.Params().Get("project")
	zip := app.GetCompressPath(r.Workspace.Path())
	ctx.SendFile(zip, pj+"-"+path.Base(zip))
}
