package app

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/cmd/docker"
)

type ServiceRoutes struct {
	*Route
}

func (r *ServiceRoutes) ListService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	out, err := docker.StackServices(stack)
	if err != nil {
		out = append(out, models.DockerStackServices{Id: err.Error()})
	}
	ctx.ViewData("out", out)
	ctx.ViewData("stack", stack)
	ctx.View("service.html")
}

func (r *ServiceRoutes) PsService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	out, err := docker.ServicePs(serviceId)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.ViewData("out", out)
	ctx.View("ps.html")
}

func (r *ServiceRoutes) RemoveService(ctx iris.Context) {
	service := ctx.Params().Get("service")
	_, _, err := docker.ServiceRm(service)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

func (r *ServiceRoutes) InspectService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	_, out, err := docker.ServiceInspect(serviceId)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}

func (r *ServiceRoutes) UpdateService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	image := ctx.FormValue("image")
	_, out, err := docker.ServiceUpdate(serviceId, "--image", image)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}

func (r *ServiceRoutes) LogsService(ctx iris.Context) {
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
