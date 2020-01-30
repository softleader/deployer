package docker

import (
	"fmt"
	"github.com/softleader/dockerctl/pkg/dockerd"
	"testing"
)

func TestParallelOverNodes(t *testing.T) {
	n := []dockerd.Node{{
		Name: "a",
		Addr: "",
	}, {
		Name: "b",
		Addr: "",
	}}
	out, err := parallelOverNodes("hello", n, func(grep string, node dockerd.Node) string {
		return fmt.Sprintf("%s %s\n", grep, node.Name)
	})
	if err != nil {
		t.Error(err)
	}
	if len(out) <= 0 {
		t.Error("out should have something")
	}
}
