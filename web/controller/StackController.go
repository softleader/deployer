package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/services"
	"time"
	"encoding/json"
	"github.com/softleader/deployer/pipe"
	"strings"
	"strconv"
)

type StackController struct {
	mvc.C
	Service services.DeployService
}

//func (c *StackController) GetStacks() mvc.Result {
//	out, err := c.Service.GetAll()
//	if err != nil {
//		out = append(out, []string{err.Error()})
//	}
//	return mvc.View{
//		Name: "stack.html",
//		Data: out,
//	}
//}

func (c *StackController) Get() mvc.Result {
	out, err := c.Service.GetAll()
	if err != nil {
		out = append(out, []string{err.Error()})
	}

	cards := make(map[string][][]string)
	for _, line := range out {
		splited := strings.Split(line[0], "-")
		key := splited[0]
		if len(splited) > 1 {
			if publishedPort(splited[1]) { // 有 publish port 可視為有開啟 dev 模式
				key += " ( port published from " + splited[1] + " )"
			}
		}
		cards[key] = append(cards[key], line)
	}
	return mvc.View{
		Name: "card.html",
		Data: cards,
	}
}

func publishedPort(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func (c *StackController) Post() {
	d := &datamodels.Deploy{}
	c.Ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")

	c.Ctx.StreamWriter(pipe.Printf("Received deploy request: %v", string(indent)))
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
