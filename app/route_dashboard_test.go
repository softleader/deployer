package app

import (
	"testing"
	"github.com/wcharczuk/go-chart"
	"fmt"
	"os"
	"image/png"
)

func TestDashboardDrawNodesChart(t *testing.T) {
	g, err := dashboardDrawNodesChart()
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	g.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("/Users/Matt/tmp/nodes.png")
	if err != nil {
		t.Error(err)
	}

	png.Encode(out, image)
	out.Close()
}

func TestDrawServicesChart(t *testing.T) {
	g, err := dashboardDrawServicesChart()
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	g.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("/Users/Matt/tmp/services.png")
	if err != nil {
		t.Error(err)
	}

	png.Encode(out, image)
	out.Close()
}

func TestDrawProjectsChart(t *testing.T) {
	g, err := dashboardDrawProjectsChart()
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	g.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("/Users/Matt/tmp/projects.png")
	if err != nil {
		t.Error(err)
	}

	png.Encode(out, image)
	out.Close()
}
