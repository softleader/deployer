package cmd

import (
	"os"
	"path"
)

type Wd struct {
	Path string
}

func NewWd() Wd {
	wd, _ := os.Getwd()
	wd = path.Join(wd, "/go")
	return Wd{Path: wd}
}

func (wd Wd) RemoveAll() Wd {
	os.RemoveAll(wd.Path)
	return wd
}

func (wd Wd) MkdirAll() Wd {
	os.MkdirAll(wd.Path, os.ModeDir|os.ModePerm)
	return wd
}
