package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"time"
	"encoding/json"
	"github.com/softleader/deployer/pipe"
)

type StackController struct {
	mvc.C
	Service services.DeployService
}

func (c *StackController) Get() mvc.Result {
	out, err := c.Service.GetAll()
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	return mvc.View{
		Name: "stack.html",
		Data: out,
	}
}

func (c *StackController) Post() {
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

func (c *StackController) GetRmBy(stack string) mvc.Result {
	_, err := c.Service.DeleteStack(stack)
	return mvc.Response{
		Err:  err,
		Path: "/",
	}
}
