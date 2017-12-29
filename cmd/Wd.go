package cmd

import "os"

// 當前的 working directory, 通常是 workspace/${project}
type Wd struct {
	Path string
}

func (wd *Wd) RemoveAll() {
	os.RemoveAll(wd.Path)
}

func (wd *Wd) MkdirAll() {
	os.MkdirAll(wd.Path, os.ModeDir|os.ModePerm)
}
