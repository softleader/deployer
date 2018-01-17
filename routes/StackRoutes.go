package routes

import (
	"github.com/kataras/iris"
	"strings"
	"time"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/services"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/pipe"
	"encoding/json"
)

type StackRoutes struct {
	args models.Args
	ds   services.DeployService
}

func NewStackRoutes(args models.Args, ds services.DeployService) *StackRoutes {
	return &StackRoutes{
		args: args,
		ds:   ds,
	}
}

func (r *StackRoutes) ListStack(ctx iris.Context) {
	out, err := r.ds.GetAll()
	if err != nil {
		out = append(out, []string{err.Error()})
	}

	stacks := make(map[string][][]string)
	for _, line := range out {
		splited := strings.Split(line[0], "-")
		key := splited[0]
		if len(splited) > 1 {
			if publishedPort(splited[1]) { // 有 publish port 可視為有開啟 dev 模式
				key = strings.Join(splited[:2], "-")
			}
		}
		_, out, _ := r.ds.DockerService.GetCreatedTimeOfFirstServiceInStack(line[0])
		out = strings.TrimSuffix(out, "\n")
		t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", out)
		line = append(line, uptime(t))
		stacks[key] = append(stacks[key], line)
	}
	ctx.ViewData("stacks", stacks)
	ctx.View("stack.html")
}

func (r *StackRoutes) DeployStack(ctx iris.Context) {
	d := &models.Deploy{}
	ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")

	ctx.StreamWriter(pipe.Printf("Received deploy request: %v", string(indent)))
	err := r.ds.Deploy(&ctx, *d)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.StreamWriter(pipe.Printf("Resolving in %v, done.", time.Since(start)))
}

func (r *StackRoutes) RemoveStack(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	_, err := r.ds.DeleteStack(stack)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.Redirect("/")
}

func uptime(t time.Time) string {
	d := time.Since(t)
	return fmt.Sprintf("up %.2f hours", d.Hours())
}

func publishedPort(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}
