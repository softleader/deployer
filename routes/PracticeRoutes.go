package routes

import (
	"github.com/softleader/deployer/models"
	"github.com/kataras/iris"
	"strings"
	"github.com/softleader/deployer/app"
)

type PracticeRoutes struct {
	app.Workspace
}

func (r *PracticeRoutes) BestPractices(ctx iris.Context) {
	out, err := models.ReadPractices(r.Workspace.Path())
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices.html")
}

func (r *PracticeRoutes) MarkdownEditor(ctx iris.Context) {
	out, err := models.ReadPractices(r.Workspace.Path())
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices-mde.html")
}

func (r *PracticeRoutes) SaveMarkdown(ctx iris.Context) {
	c := ctx.PostValue("content")
	c = strings.TrimSpace(c)
	if len(c) > 0 {
		err := models.SavePractices(r.Workspace.Path(), c)
		if err != nil {
			ctx.Application().Logger().Warn(err.Error())
			ctx.WriteString(err.Error())
		}
	}
	ctx.Redirect("/best-practices")
}
