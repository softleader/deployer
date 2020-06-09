package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/cmd/docker"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/slack"
	"strconv"
	"strings"
)

func FilterService(ctx iris.Context) {
	params := ctx.URLParams()
	_, out, err := docker.ServiceFilter(params, false)
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
	image := ctx.FormValue("image")
	replicas := ctx.FormValue("replicas")
	if i, err := strconv.Atoi(replicas); err != nil || i < 0 {
		replicas = ""
	}
	if image == "" && replicas == "" {
		writeOut(ctx, "Requires at least one of 'image' or 'replicas' parameter", 400)
		return
	}
	serviceId := ctx.Params().Get("serviceId")
	if filter := ctx.FormValue("filter"); filter != "" {
		params := make(map[string]string)
		params[filter] = serviceId
		arg, ids, err := findServiceIdByLabel(params)
		if err != nil {
			writeOut(ctx, err.Error(), 400)
			return
		}
		if len(ids) == 0 {
			writeOut(ctx, fmt.Sprintf("No service found for: %s", arg), 404)
			return
		}
		if len(ids) > 1 {
			writeOut(ctx, fmt.Sprintf("No unique service found for: %s", arg), 400)
			return
		}
		serviceId = ids[0]
	}
	if _, found := ctx.FormValues()["skip-slack"]; !found {
		err := slack.Post(Ws.Config, serviceId, image)
		if err != nil {
			fmt.Println(err)
		}
	}
	var args []string
	if image != "" {
		args = append(args, "--image", image)
	}
	if replicas != "" {
		args = append(args, "--replicas", replicas)
	}
	_, out, err := docker.ServiceUpdate(serviceId, args...)
	if err != nil {
		out += err.Error()
		writeOut(ctx, out, 400)
		return
	}
	writeOut(ctx, out, 202)
}

func findServiceIdByLabel(params map[string]string) (arg string, ids []string, err error) {
	var out string
	if arg, out, err = docker.ServiceFilter(params, true); err != nil {
		return
	}
	ids = deleteEmpty(strings.Split(out, fmt.Sprintln()))
	return
}

func writeOut(ctx iris.Context, out string, status int) {
	if ctx.Method() != "GET" {
		ctx.StatusCode(status)
		ctx.Text(out)
		return
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
