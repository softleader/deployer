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
	WarningStyle = chart.Style{
		Show:      true,
		FillColor: chart.ColorAlternateYellow,
	}
	ErrorStyle = chart.Style{
		Show:      true,
		FillColor: chart.ColorRed,
		FontColor: chart.ColorWhite,
	}
	SuccessStyle = chart.Style{
		Show:      true,
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
	m := make(map[string]chart.Value)
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			node := strings.Split(line, ";")
			host := node[0]
			status := node[1]
			if status == "Down" {
				v := m[status].Value + 1
				l := m[status].Label + host + ", "
				m[status] = chart.Value{
					Value: v,
					Style: ErrorStyle,
					Label: l,
				}
			} else {
				v := m[status].Value + 1
				m[status] = chart.Value{Value: v, Label: fmt.Sprintf("%s (%v)", status, v)}
			}
		}
	}
	down := m["Down"]
	if down.Value > 0 {
		m["Down"] = chart.Value{
			Value: down.Value,
			Style: ErrorStyle,
			Label: fmt.Sprintf("Down (%s)", strings.TrimSuffix(down.Label, ", ")),
		}
	}
	var values []chart.Value
	for _, v := range m {
		values = append(values, v)
	}
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	flush(ctx, pie)
}

func (r *DashboardRoutes) Projects(ctx iris.Context) {
	out, err := r.DockerStack.Ls()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}

	projects := make(map[string][]string)
	for _, stack := range out {
		p := strings.Split(stack[0], "-")[0]
		for _, name := range stack {
			projects[p] = append(projects[p], name)
		}
	}

	var bars []chart.StackedBar

	for pj, stacks := range projects {
		var services [][]string
		for _, stack := range stacks {
			svcs, _ := r.DockerStack.Services(stack)
			for _, svc := range svcs {
				services = append(services, svc)
			}
		}

		label, values := toStackedBarValues(services)
		bars = append(bars, chart.StackedBar{
			Name:   fmt.Sprintf("%s %s", pj, label),
			Values: values,
		})
	}

	if len(bars) == 0 {
		bars = append(bars, chart.StackedBar{
			Name: "No project found",
			Values: []chart.Value{{Value: 1, Style: WarningStyle},
			},
		})
	}

	sbc := chart.StackedBarChart{
		Height: 512,
		XAxis:  chart.StyleShow(),
		Bars:   bars,
	}
	flush(ctx, sbc)
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
					Style: SuccessStyle,
					Label: fmt.Sprintf("Healthy (%.0f)", v),
				}
			} else {
				v := m["Unhealthy"].Value + 1
				m["Unhealthy"] = chart.Value{
					Value: m["Unhealthy"].Value + 1,
					Style: ErrorStyle,
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
			Style: WarningStyle,
		})
	}
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	flush(ctx, pie)
}

func toStackedBarValues(svcs [][]string) (label string, values chart.Values) {
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
				Style: SuccessStyle,
				Label: fmt.Sprintf("Healthy (%.0f)", v),
			}
			healthy++
		} else {
			v := m["Unhealthy"].Value + 1
			m["Unhealthy"] = chart.Value{
				Value: m["Unhealthy"].Value + 1,
				Style: ErrorStyle,
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
