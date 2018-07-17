package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"path"
	"strconv"
)

type DeployRoutes struct {
	Routes
}

func (r *DeployRoutes) DeployPage(ctx iris.Context) {
	ctx.ViewData("workspace", r.Args.Ws)

	h := ctx.Params().Get("history")
	dft, err := prepareDefaultValue(r.Workspace, h)
	if err != nil {
		ctx.ViewData("err", err)
	}
	ctx.ViewData("dft", dft)
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.View("deploy.html")
}

func prepareDefaultValue(ws *Workspace, h string) (d models.Deploy, err error) {
	d = ws.Config.Deploy
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
	zip := GetCompressPath(path.Join(r.Workspace.Path(), pj))
	ctx.SendFile(zip, pj+"-"+path.Base(zip))
}
