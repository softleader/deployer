package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
)

func GetHistory(ctx iris.Context) {
	out, err := models.GetHistory(Ws.Path())
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("history.html")
}

func RemoveHistory(ctx iris.Context) {
	index, err := ctx.Params().GetInt("idx")
	if err != nil {
		ctx.ViewData("err", err)
		return
	}
	h, err := models.GetHistory(Ws.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
	h.Delete(index)
	h.SaveTo(Ws.Path())
	if err != nil {
		ctx.ViewData("err", err)
	}
}
