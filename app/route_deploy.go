package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"strconv"
	"path"
)

func DeployPage(ctx iris.Context) {
	ctx.ViewData("workspace", Args.Ws)

	h := ctx.Params().Get("history")
	dft, err := prepareDefaultValue(h)
	if err != nil {
		ctx.ViewData("err", err)
	}
	ctx.ViewData("dft", dft)
	ctx.View("deploy.html")
}

func prepareDefaultValue(h string) (d models.Deploy, err error) {
	d = Ws.Config.Deploy
	if h != "" {
		i, err := strconv.Atoi(h)
		if err == nil && i >= 0 {
			histories, err := models.GetHistory(Ws.Path())
			if err == nil && i < len(histories) {
				d = histories[i]
			}
		}
	}
	return d, err
}

func DownloadYAML(ctx iris.Context) {
	pj := ctx.Params().Get("project")
	zip := GetCompressPath(path.Join(Ws.Path(), pj))
	ctx.SendFile(zip, pj+"-"+path.Base(zip))
}
