package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
)

// workspace, 一個 golang app 會有唯一的一個 workspace
type Workspace struct {
	path string
}

func NewWorkspace(dir string) *Workspace {
	if dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/workspace")
	}
	wd := Workspace{path: dir}
	fmt.Printf("Setting up workspace to '%v'\n", dir)

	stat, err := os.Stat(wd.Path())

	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			stat, err = os.Stat(wd.Path())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		} else {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	if !stat.IsDir() {
		log.Fatal(fmt.Sprintf("Workspace requires a dictionary: %v", wd.Path()))
		os.Exit(1)
	}

	_, err = os.Open(wd.Path())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &wd
}

func (ws *Workspace) GetWd(cleanUp bool, project string) *WorkDir {
	return NewWorkDir(cleanUp, path.Join(ws.Path(), project))
}

func (ws *Workspace) Path() string {
	return ws.path
}
