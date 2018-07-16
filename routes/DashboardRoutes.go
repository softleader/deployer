package routes

import (
	"github.com/kataras/iris"
	"github.com/softleader/deployer/app"
	"github.com/softleader/deployer/cmd"
	"github.com/wcharczuk/go-chart"
	"strings"
	"strconv"
	"fmt"
	"io"
)

var (
	WARNING_STYLE = chart.Style{
		FillColor: chart.ColorAlternateYellow,
	}
	ERROR_STYLE = chart.Style{
		FillColor: chart.ColorRed,
	}
	SUCCESS_STYLE = chart.Style{
		FillColor: chart.ColorAlternateGreen,
	}
)

type DashboardRoutes struct {
	app.Workspace
	cmd.DockerNode
	cmd.DockerStack
	cmd.DockerService
}

func (r *DashboardRoutes) DashboardPage(ctx iris.Context) {
	ctx.ViewData("navbar", r.Workspace.Config.Navbar)
	ctx.View("dashboard.html")
}

func (r *DashboardRoutes) Nodes(ctx iris.Context) {
	_, out, err := r.DockerNode.Ls()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	lines := strings.Split(out, "\n")
	m := make(map[string]float64)
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			node := strings.Split(line, ";")
			m[node[0]]++
		}
	}
	var values []chart.Value
	for k, v := range m {
		values = append(values, chart.Value{Value: v, Label: fmt.Sprintf("%s (%.0f)", k, v)})
	}
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	flush(ctx, pie)
}

func (r *DashboardRoutes) Stacks(ctx iris.Context) {
	out, err := r.DockerStack.Ls()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}

	var bars []chart.StackedBar

	for _, stack := range out {
		svcs, _ := r.DockerStack.Services(stack[0])
		label, values := toStackValues(svcs)
		name := fmt.Sprintf("%s %s", stack[0], label)
		bars = append(bars, chart.StackedBar{
			Name:   escape(name),
			Values: values,
		})
	}

	sbc := chart.StackedBarChart{
		Height: 512,
		XAxis:  chart.StyleShow(),
		Bars:   bars,
	}
	flush(ctx, sbc)
}

// 有些字元會造成 chart 無法 render..
func escape(name string) (escaped string) {
	escaped = strings.Replace(name, "-", " ", -1)
	escaped = strings.Replace(escaped, "_", " ", -1)
	return
}

func (r *DashboardRoutes) Services(ctx iris.Context) {
	_, out, err := r.DockerService.Ls()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	lines := strings.Split(out, "\n")
	m := make(map[string]chart.Value)
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			svc := strings.Split(line, ";")
			replicas := strings.Split(svc[0], "/")
			up, _ := strconv.Atoi(replicas[0])
			total, _ := strconv.Atoi(replicas[1])
			if up == total {
				v := m["Healthy"].Value + 1
				m["Healthy"] = chart.Value{
					Value: v,
					Style: SUCCESS_STYLE,
					Label: fmt.Sprintf("Healthy (%.0f)", v),
				}
			} else {
				v := m["Unhealthy"].Value + 1
				m["Unhealthy"] = chart.Value{
					Value: m["Unhealthy"].Value + 1,
					Style: ERROR_STYLE,
					Label: fmt.Sprintf("Unhealthy (%.0f)", v),
				}
			}
		}
	}
	var values []chart.Value
	for _, v := range m {
		values = append(values, v)
	}
	if len(values) == 0 {
		values = append(values, chart.Value{
			Value: 1,
			Label: "No service found",
			Style: WARNING_STYLE,
		})
	}
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	flush(ctx, pie)
}

func toStackValues(svcs [][]string) (label string, values chart.Values) {
	var healthy, unhealthy int
	m := make(map[string]chart.Value)
	for _, svc := range svcs {
		replicas := strings.Split(svc[3], "/")
		up, _ := strconv.Atoi(replicas[0])
		total, _ := strconv.Atoi(replicas[1])
		if up == total {
			v := m["Healthy"].Value + 1
			m["Healthy"] = chart.Value{
				Value: v,
				Style: SUCCESS_STYLE,
				Label: fmt.Sprintf("Healthy (%.0f)", v),
			}
			healthy++
		} else {
			v := m["Unhealthy"].Value + 1
			m["Unhealthy"] = chart.Value{
				Value: m["Unhealthy"].Value + 1,
				Style: ERROR_STYLE,
				Label: fmt.Sprintf("Unhealthy (%.0f)", v),
			}
			unhealthy++
		}
	}
	for _, v := range m {
		values = append(values, v)
	}
	label = fmt.Sprintf("(%d/%d)", healthy, healthy+unhealthy)
	return
}

func flush(ctx iris.Context, r Render) {
	ctx.Header("Content-Type", "image/png")
	err := r.Render(chart.PNG, ctx.ResponseWriter())
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

type Render interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}
