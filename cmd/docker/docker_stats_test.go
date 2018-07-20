package docker

import (
	"testing"
	"fmt"
)

func TestParallelOverNodes(t *testing.T) {
	n := []string{"a", "b"}
	out, err := parallelOverNodes("hello", n, func(grep string, host string) string {
		return fmt.Sprintf("%s %s\n", grep, host)
	})
	if err != nil {
		t.Error(err)
	}
	if len(out) <= 0 {
		t.Error("out should have something")
	}
}
