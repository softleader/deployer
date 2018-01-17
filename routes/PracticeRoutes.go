package routes

import (
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/services"
	"github.com/kataras/iris"
)

type PracticeRoutes struct {
	args models.Args
	ps   services.PracticeService
}

func NewPracticeRoutes(args models.Args, ps services.PracticeService) *PracticeRoutes {
	return &PracticeRoutes{
		args: args,
		ps:   ps,
	}
}

func (r *PracticeRoutes) BestPractices(ctx iris.Context) {
	out, err := r.ps.Get()
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices.html")
}

func (r *PracticeRoutes) MarkdownEditor(ctx iris.Context) {
	out, err := r.ps.Get()
	ctx.ViewData("err", err)
	ctx.ViewData("out", out)
	ctx.View("best-practices-mde.html")
}

func (r *PracticeRoutes) SaveMarkdown(ctx iris.Context) {
	c := ctx.PostValue("content");
	err := r.ps.Save(c)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.Redirect("/best-practices")
}
