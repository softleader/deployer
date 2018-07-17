package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
)

type HistoryRoutes struct {
	Routes
}

func (r *HistoryRoutes) GetHistory(ctx iris.Context) {
	out, err := models.GetHistory(r.Workspace.Path())
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("history.html")
}

func (r *HistoryRoutes) RemoveHistory(ctx iris.Context) {
	index, err := ctx.Params().GetInt("idx")
	if err != nil {
		ctx.ViewData("err", err)
		return
	}
	h, err := models.GetHistory(r.Workspace.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
	h.Delete(index)
	h.SaveTo(r.Workspace.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
}
