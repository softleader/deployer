package routes

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/app"
	"path"
	"strconv"
)

type DeployRoutes struct {
	app.Args
	app.Workspace
}

func (r *DeployRoutes) DeployPage(ctx iris.Context) {
	ctx.ViewData("workspace", r.Args.Ws)

	h := ctx.Params().Get("history")
	dft, err := prepareDefaultValue(r.Workspace, h)
	if err != nil {
		ctx.ViewData("err", err)
	}
	ctx.ViewData("dft", dft)
	ctx.View("deploy.html")
}

func prepareDefaultValue(ws app.Workspace, h string) (d models.Deploy, err error) {
	d = *models.NewDefaultDeploy(ws.Path())
	if h != "" {
		i, err := strconv.Atoi(h)
		if err == nil && i >= 0 {
			histories, err := models.GetHistory(ws.Path())
			if err == nil && i < len(histories) {
				d = histories[i]
			}
		}
	}
	return d, err
}

func (r *DeployRoutes) DownloadYAML(ctx iris.Context) {
	pj := ctx.Params().Get("project")
	zip := app.GetCompressPath(path.Join(r.Workspace.Path(), pj))
	ctx.SendFile(zip, pj+"-"+path.Base(zip))
}
