package routes

import (
	"testing"
	"github.com/softleader/deployer/cmd"
	"github.com/wcharczuk/go-chart"
	"fmt"
	"os"
	"image/png"
	"github.com/softleader/deployer/app"
)

func TestDrawNodesChart(t *testing.T) {
	cmd := cmd.NewDockerNode()
	cr, err := drawNodesChart(*cmd)
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	cr.Render(chart.PNG, collector)

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
	cmd := cmd.NewDockerService()
	cr, err := drawServicesChart(*cmd)
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	cr.Render(chart.PNG, collector)

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
	cmd := cmd.NewDockerStack(app.Registry{})
	cr, err := drawProjectsChart(*cmd)
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	cr.Render(chart.PNG, collector)

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

func TestDrawStatsChart(t *testing.T) {
	cmd := cmd.NewDockerStats()
	cr, err := drawStatsChart(*cmd, "test")
	if err != nil {
		t.Error(err)
	}
	collector := &chart.ImageWriter{}
	cr.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	out, err := os.Create("/Users/Matt/tmp/stats.png")
	if err != nil {
		t.Error(err)
	}

	png.Encode(out, image)
	out.Close()
}
