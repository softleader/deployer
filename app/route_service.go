package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/cmd/docker"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/slack"
	"strings"
	"time"
)

func FilterService(ctx iris.Context) {
	params := ctx.URLParams()
	_, out, err := docker.ServiceFilter(params)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.WriteString(fmt.Sprintf("[%s]", strings.Join(deleteEmpty(strings.Split(out, fmt.Sprintln())), ",")))
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func ListService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	out, err := docker.StackServices(stack)
	if err != nil {
		out = append(out, models.DockerStackServices{Id: err.Error()})
	}
	ctx.ViewData("out", out)
	ctx.ViewData("stack", stack)
	ctx.View("service.html")
}

func PsService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	out, err := docker.ServicePs(serviceId)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.ViewData("out", out)
	ctx.View("ps.html")
}

func RemoveService(ctx iris.Context) {
	service := ctx.Params().Get("service")
	_, _, err := docker.ServiceRm(service)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

func InspectService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	_, out, err := docker.ServiceInspect(serviceId)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}

func UpdateService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	image := ctx.FormValue("image")
	if hookSlack, _ := ctx.Params().GetBool("slack"); hookSlack {
		title := ctx.Params().Get("title")
		titleLink := ctx.Params().Get("title_link")
		authorLink := ctx.Params().Get("author_name")
		authorName := ctx.Params().Get("author_link")
		authorIcon := ctx.Params().Get("author_icon")
		ts, err := ctx.Params().GetInt64("ts")
		if err != nil {
			ts = time.Now().Unix()
		}
		slack.Post(Ws.Config.SlackAPI, image, title, titleLink, authorLink, authorName, authorIcon, ts)
	}
	_, out, err := docker.ServiceUpdate(serviceId, "--image", image)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}

func LogsService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	tail, err := ctx.Params().GetInt("tail")
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	if tail <= 0 {
		tail = 300
	}
	_, out, err := docker.ServiceLogs(serviceId, tail)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}
