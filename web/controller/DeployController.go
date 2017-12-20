package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"fmt"
	"time"
	"encoding/json"
	"github.com/kataras/iris"
	"io"
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

	c.Ctx.StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, "Receiving deploy request: %v", string(indent))
		return false
	})
	err := c.Service.Deploy(&c.Ctx, *d)
	if err != nil {
		c.Ctx.Application().Logger().Warn(err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, "Resolving in '%v', done.", time.Since(start))
		return false
	})
	return
}

func (c *DeployController) DeleteBy(stack string) (string, int) {
	out, err := c.Service.Delete(stack)
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}
