package cmd

import (
	"os"
	"path"
)

type Wd struct {
	Path string
}

func NewWd() Wd {
	pwd, _ := os.Getwd()
	pwd = path.Join(pwd, "/go")
	wd := Wd{Path: pwd}
	if _, err := os.Stat(pwd); os.IsNotExist(err) {
		wd.MkdirAll()
	}
	return wd
}

func (wd Wd) RemoveAll() Wd {
	os.RemoveAll(wd.Path)
	return wd
}

func (wd Wd) MkdirAll() Wd {
	os.MkdirAll(wd.Path, os.ModeDir|os.ModePerm)
	return wd
}
