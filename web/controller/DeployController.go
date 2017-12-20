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

func (c *DeployController) Get() string {
	out, err := c.Service.GetAll()
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		return err.Error()
	}
	return out
}

func (c *DeployController) GetBy(stack string) string {
	out, err := c.Service.GetServices(stack)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		return err.Error()
	}
	return out
}

func (c *DeployController) Post() string {
	d := &datamodels.Deploy{}
	c.Ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")
	var buf bytes.Buffer
	msg := fmt.Sprintf("\n[%v] Receiving %v\n", start.Format(time.Stamp), string(indent))
	fmt.Print(msg)
	buf.WriteString(msg)
	out, err := c.Service.Deploy(*d)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		buf.WriteString(err.Error())
	}
	buf.WriteString(out)
	msg = fmt.Sprintf("[%v] Resolving in '%v', done.\n", time.Now().Format(time.Stamp), time.Since(start))
	fmt.Print(msg)
	buf.WriteString("\n" + msg)
	return buf.String()
}

func (c *DeployController) DeleteBy(stack string) string {
	out, err := c.Service.Delete(stack)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		return err.Error()
	}
	return out
}
