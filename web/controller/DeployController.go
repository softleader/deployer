package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	"github.com/kataras/iris"
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

func (c *DeployController) Post() (string, int) {
	d := &datamodels.Deploy{}
	c.Ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")
	var buf bytes.Buffer
	msg := fmt.Sprintf("Receiving deploy request: %v", string(indent))
	c.Ctx.Application().Logger().Info(msg)
	buf.WriteString(msg)
	out, err := c.Service.Deploy(*d)
	if err != nil {
		c.Ctx.Application().Logger().Warn(err.Error())
		buf.WriteString(err.Error())
		return buf.String(), iris.StatusInternalServerError
	}
	buf.WriteString(out)
	msg = fmt.Sprintf("Resolving in '%v', done.", time.Since(start))
	c.Ctx.Application().Logger().Info(msg)
	buf.WriteString("\n" + msg)
	return buf.String(), iris.StatusOK
}

func (c *DeployController) DeleteBy(stack string) (string, int) {
	out, err := c.Service.Delete(stack)
	if err != nil {
		return err.Error(), iris.StatusInternalServerError
	}
	return out, iris.StatusOK
}
