package models

import (
	"testing"
	"sort"
)

func TestSaveOrUpdate(t *testing.T) {

	h := Histories{
		{
			Project: "c",
			Yaml:    "hello",
		},
		{
			Project: "b",
			Yaml:    "",
		},
		{
			Project: "a",
			Yaml:    "world",
		},
	}

	d := Deploy{
		Project: "b",
		Yaml:    "my-yaml",
	}

	h.Push(&d)

	errorIfNotEquals(t, "length", 3, len(h))
	errorIfNotEquals(t, "[0].Project", "c", h[0].Project)
	errorIfNotEquals(t, "[0].Yaml", "hello", h[0].Yaml)
	errorIfNotEquals(t, "[1].Project", "b", h[1].Project)
	errorIfNotEquals(t, "[1].Yaml", "my-yaml", h[1].Yaml)
	errorIfNotEquals(t, "[2].Project", "a", h[2].Project)
	errorIfNotEquals(t, "[2].Yaml", "world", h[2].Yaml)

	sort.Sort(h)

	errorIfNotEquals(t, "[0].Project", "a", h[0].Project)
	errorIfNotEquals(t, "[0].Yaml", "world", h[0].Yaml)
	errorIfNotEquals(t, "[1].Project", "b", h[1].Project)
	errorIfNotEquals(t, "[1].Yaml", "my-yaml", h[1].Yaml)
	errorIfNotEquals(t, "[2].Project", "c", h[2].Project)
	errorIfNotEquals(t, "[2].Yaml", "hello", h[2].Yaml)

}

func errorIfNotEquals(t *testing.T, n string, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Error("expected", n, "equals to", expected, ", but actual is", actual)
	}
}
