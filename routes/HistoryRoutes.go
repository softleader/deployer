package routes

import (
	"github.com/softleader/deployer/app"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
)

type HistoryRoutes struct {
	app.Workspace
}

func (r *HistoryRoutes) GetHistories(ctx iris.Context) {
	out, err := models.GetHistories(r.Workspace.Path())
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("histories.html")
}

func (r *HistoryRoutes) RemoveHistory(ctx iris.Context) {
	index, err := ctx.Params().GetInt("index")
	if err != nil {
		ctx.ViewData("err", err)
		return
	}
	h, err := models.GetHistories(r.Workspace.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
	h.Delete(index)
	h.SaveTo(r.Workspace.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
}
