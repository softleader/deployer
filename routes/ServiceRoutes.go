package routes

import (
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/services"
	"github.com/kataras/iris"
)

type ServiceRoutes struct {
	args models.Args
	ds   services.DeployService
}

func NewServiceRoutes(args models.Args, ds services.DeployService) *ServiceRoutes {
	return &ServiceRoutes{
		args: args,
		ds:   ds,
	}
}

func (r *ServiceRoutes) ListService(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	out, err := r.ds.GetServices(stack)
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	ctx.ViewData("out", out)
	ctx.ViewData("stack", stack)
	ctx.View("service.html")
}

func (r *ServiceRoutes) PsService(ctx iris.Context) {
	serviceId := ctx.Params().Get("serviceId")
	out, err := r.ds.Ps(serviceId)
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
	_, err := r.ds.DeleteService(service)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.Redirect("/services/" + stack)
}
