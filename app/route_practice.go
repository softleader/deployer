package app

import (
	"github.com/softleader/deployer/models"
	"github.com/kataras/iris"
	"strings"
)

func BestPractices(ctx iris.Context) {
	out, err := models.ReadPractices(Ws.Path())
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices.html")
}

func MarkdownEditor(ctx iris.Context) {
	out, err := models.ReadPractices(Ws.Path())
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices-mde.html")
}

func SaveMarkdown(ctx iris.Context) {
	c := ctx.PostValue("content")
	c = strings.TrimSpace(c)
	if len(c) > 0 {
		err := models.SavePractices(Ws.Path(), c)
		if err != nil {
			ctx.Application().Logger().Warn(err.Error())
			ctx.WriteString(err.Error())
		}
	}
	ctx.Redirect("/best-practices")
}
