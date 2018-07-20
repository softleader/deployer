package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"strings"
	"github.com/softleader/deployer/cmd/docker"
)

type StatsRoutes struct {
	*Route
}

func (r *StatsRoutes) GetStats(ctx iris.Context) {
	out, err := docker.StackLs()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		return
	}

	projects := make(map[string]struct{})
	for _, stack := range out {
		projects[models.NewStackName(stack.Name).Project] = struct{}{}
	}
	ctx.ViewData("projects", projects)

	grep := ctx.FormValue("g")
	if g := strings.TrimSpace(grep); g != "" {
		out, err := docker.StatsNoStream(g)
		ctx.ViewData("err", err)
		ctx.ViewData("out", out)
	}

	ctx.View("stats.html")
}
