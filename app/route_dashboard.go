package app

import (
	"github.com/kataras/iris"
	"github.com/wcharczuk/go-chart"
	"strings"
	"fmt"
	"io"
	"github.com/softleader/deployer/models"
	"github.com/softleader/deployer/cmd/docker"
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

func DashboardPage(ctx iris.Context) {
	ctx.View("dashboard.html")
}

func DashboardNodes(ctx iris.Context) {
	g, err := dashboardDrawNodesChart()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, g)
}

func dashboardDrawNodesChart() (g Graph, err error) {
	out, err := docker.NodeLs()
	if err != nil {
		return
	}
	m := make(map[string]chart.Value)
	for _, node := range out {
		if node.Status == "Down" {
			v := m[node.Status].Value + 1
			l := m[node.Status].Label + node.Hostname + ", "
			m[node.Status] = chart.Value{
				Value: v,
				Style: ErrorStyle,
				Label: l,
			}
		} else {
			v := m[node.Status].Value + 1
			m[node.Status] = chart.Value{Value: v, Label: fmt.Sprintf("%s (%v)", node.Status, v)}
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
	g = chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	return
}

func DashboardProjects(ctx iris.Context) {
	g, err := dashboardDrawProjectsChart()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, g)
}

func dashboardDrawProjectsChart() (r Graph, err error) {
	out, err := docker.StackLs()
	if err != nil {
		return
	}

	projects := make(map[string][]string)
	for _, stack := range out {
		p := models.NewStackName(stack.Name).Project
		projects[p] = append(projects[p], stack.Name)
	}

	var bars []chart.StackedBar

	for pj, stacks := range projects {
		var services []models.DockerStackServices
		for _, s := range stacks {
			svcs, _ := docker.StackServices(s)
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

	r = chart.StackedBarChart{
		Height: 512,
		XAxis:  chart.StyleShow(),
		Bars:   bars,
	}
	return
}

func DashboardServices(ctx iris.Context) {
	g, err := dashboardDrawServicesChart()
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, g)
}

func dashboardDrawServicesChart() (g Graph, err error) {
	out, err := docker.ServiceLs()
	if err != nil {
		return
	}
	m := make(map[string]chart.Value)
	for _, svc := range out {
		if svc.Rtoi().LooksGood() {
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
	g = chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	return
}

func toStackedBarValues(svcs []models.DockerStackServices) (label string, values chart.Values) {
	var healthy, unhealthy int
	m := make(map[string]chart.Value)
	for _, svc := range svcs {
		if svc.Rtoi().LooksGood() {
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

func flush(ctx iris.Context, g Graph) {
	ctx.Header("Content-Type", "image/png")
	err := g.Render(chart.PNG, ctx.ResponseWriter())
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

type Graph interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}
