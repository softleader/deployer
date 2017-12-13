package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
)

type DeploymentController struct {
	mvc.C
	Service services.DeploymentService
}

func (c *DeploymentController) Get() string {
	return c.Service.GetAll()
}

func (c *DeploymentController) PostBy(stack string) string {
	d := &datamodels.Deployment{}
	c.Ctx.ReadJSON(d)
	return c.Service.Deploy(stack, *d)
}

func (c *DeploymentController) DeleteBy(stack string) string {
	return c.Service.Delete(stack)
}
