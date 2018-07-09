package routes

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/app"
)

type ServiceRoutes struct {
	Workspace app.Workspace
	cmd.DockerStack
	cmd.DockerService
}

func (r *ServiceRoutes) ListService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	out, err := r.DockerStack.Services(stack)
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
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
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.ViewData("out", out)
	ctx.View("ps.html")
}

func (r *ServiceRoutes) RemoveService(ctx iris.Context) {
	service := ctx.Params().Get("service")
	_, _, err := r.DockerService.Rm(service)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

func (r *ServiceRoutes) InspectService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	_, out, err := r.DockerService.Inspect(serviceId)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
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
	_, out, err := r.DockerService.Logs(serviceId, tail)
	if err != nil {
		out += err.Error()
	}
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.ViewData("out", out)
	ctx.View("pre.html")
}
