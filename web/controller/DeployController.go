package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"time"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/pipe"
)

type DeployController struct {
	mvc.C
	Service services.DeployService
}

func (c *DeployController) Get() (string, int) {
	out, err := c.Service.GetAll()
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}

func (c *DeployController) GetBy(stack string) (string, int) {
	out, err := c.Service.GetServices(stack)
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}

func (c *DeployController) GetPsBy(serviceId string) (string, int) {
	out, err := c.Service.Ps(serviceId)
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}

func (c *DeployController) Post() {
	d := &datamodels.Deploy{}
	c.Ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")

	c.Ctx.StreamWriter(pipe.Printf("Receiving deploy request: %v", string(indent)))
	err := c.Service.Deploy(&c.Ctx, *d)
	if err != nil {
		c.Ctx.Application().Logger().Warn(err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.StreamWriter(pipe.Printf("Resolving in %v, done.", time.Since(start)))
	return
}

func (c *DeployController) DeleteBy(stack string) (string, int) {
	out, err := c.Service.Delete(stack)
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}
