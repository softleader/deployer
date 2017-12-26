package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/services"
)

type ServiceController struct {
	mvc.C
	Service services.DeployService
}

func (c *ServiceController) GetBy(stack string) mvc.Result {
	out, err := c.Service.GetServices(stack)
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	return mvc.View{
		Name: "service.html",
		Data: map[string]interface{}{
			"out":   out,
			"stack": stack,
		},
	}
}

func (c *ServiceController) GetPsBy(serviceId string) mvc.Result {
	out, err := c.Service.Ps(serviceId)
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	return mvc.View{
		Name: "ps.html",
		Data: out,
	}
}

func (c *ServiceController) GetRmBy(stack string, service string) mvc.Result {
	_, err := c.Service.DeleteService(service)
	return mvc.Response{
		Err:  err,
		Path: "/services/" + stack,
	}
}
