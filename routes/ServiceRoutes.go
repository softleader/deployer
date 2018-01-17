package routes

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/cmd"
)

type ServiceRoutes struct {
	cmd.DockerStack
	cmd.DockerService
}

func (r *ServiceRoutes) ListService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	out, err := r.DockerStack.Services(stack)
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	ctx.ViewData("out", out)
	ctx.ViewData("stack", stack)
	ctx.View("service.html")
}

func (r *ServiceRoutes) PsService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	out, err := r.DockerService.Ps(serviceId)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.ViewData("out", out)
	ctx.View("ps.html")
}

func (r *ServiceRoutes) RemoveService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	service := ctx.Params().Get("service")
	_, _, err := r.DockerService.Rm(service)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.Redirect("/services/" + stack)
}
