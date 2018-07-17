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
	"github.com/softleader/deployer/models"
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
	cr, err := drawNodesChart(r.DockerNode)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, cr)
}

func drawNodesChart(node cmd.DockerNode) (cr ChartRenderer, err error) {
	out, err := node.Ls()
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
	cr = chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	return
}

func (r *DashboardRoutes) Projects(ctx iris.Context) {
	cr, err := drawProjectsChart(r.DockerStack)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, cr)
}

func drawProjectsChart(stack cmd.DockerStack) (r ChartRenderer, err error) {
	out, err := stack.Ls()
	if err != nil {
		return
	}

	projects := make(map[string][]string)
	for _, stack := range out {
		p := strings.Split(stack.Name, "-")[0]
		projects[p] = append(projects[p], stack.Name)
	}

	var bars []chart.StackedBar

	for pj, stacks := range projects {
		var services []models.DockerStackServices
		for _, s := range stacks {
			svcs, _ := stack.Services(s)
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

func (r *DashboardRoutes) Services(ctx iris.Context) {
	cr, err := drawServicesChart(r.DockerService)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	flush(ctx, cr)
}

func drawServicesChart(service cmd.DockerService) (cr ChartRenderer, err error) {
	_, out, err := service.Ls()
	if err != nil {
		return
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
	cr = chart.PieChart{
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
		replicas := strings.Split(svc.Replicas, "/")
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

func drawStatsChart(stats cmd.DockerStats, project string) (cr ChartRenderer, err error) {
	_, out, err := stats.NoStream(project)
	if err != nil {
		return
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
	cr = chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}
	return
}

func flush(ctx iris.Context, cr ChartRenderer) {
	ctx.Header("Content-Type", "image/png")
	err := cr.Render(chart.PNG, ctx.ResponseWriter())
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
}

type ChartRenderer interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}
