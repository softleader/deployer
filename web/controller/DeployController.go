package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"fmt"
	"time"
	"encoding/json"
	"bytes"
)

type DeployController struct {
	mvc.C
	Service services.DeployService
}

func (c *DeployController) Get() string {
	return c.Service.GetAll()
}

func (c *DeployController) GetBy(stack string) string {
	return c.Service.GetServices(stack)
}

func (c *DeployController) Post() string {
	d := &datamodels.Deploy{}
	c.Ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")
	var resp bytes.Buffer
	msg := fmt.Sprintf("\n[%v] Receiving %v\n", start.Format(time.Stamp), string(indent))
	fmt.Print(msg)
	resp.WriteString(msg)
	resp.WriteString(c.Service.Deploy(*d))
	msg = fmt.Sprintf("[%v] Resolving in '%v', done.\n", time.Now().Format(time.Stamp), time.Since(start))
	fmt.Print(msg)
	resp.WriteString("\n" + msg)
	return resp.String()
}

func (c *DeployController) DeleteBy(stack string) string {
	return c.Service.Delete(stack)
}
