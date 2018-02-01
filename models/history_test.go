package models

import (
	"testing"
)

func TestSaveOrUpdate(t *testing.T) {

	h := []Deploy{
		{
			Project: "a",
			Yaml:    "hello",
		},
		{
			Project: "b",
			Yaml:    "",
		},
		{
			Project: "c",
			Yaml:    "world",
		},
	}

	d := Deploy{
		Project: "b",
		Yaml:    "my-yaml",
	}

	actual := saveOrUpdate(&d, h)

	errorIfNotEquals(t, "length", 3, len(actual))
	errorIfNotEquals(t, "[0].Project", "a", actual[0].Project)
	errorIfNotEquals(t, "[0].Yaml", "hello", actual[0].Yaml)
	errorIfNotEquals(t, "[1].Project", "b", actual[1].Project)
	errorIfNotEquals(t, "[1].Yaml", "my-yaml", actual[1].Yaml)
	errorIfNotEquals(t, "[2].Project", "c", actual[2].Project)
	errorIfNotEquals(t, "[2].Yaml", "world", actual[2].Yaml)

}

func errorIfNotEquals(t *testing.T, n string, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Error("expected", n, "equals to", expected, ", but actual is", actual)
	}
}
