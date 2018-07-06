package app

import (
	"os"
	"path"
	"log"
	"fmt"
	"github.com/softleader/deployer/models"
)

// workspace, 一個 golang app 會有唯一的一個 workspace
type Workspace struct {
	path   string
	Config models.Config
}

func NewWorkspace(dir string) *Workspace {
	if dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/workspace")
	}
	ws := Workspace{path: dir}
	fmt.Printf("Setting up workspace to '%v'\n", dir)

	stat, err := os.Stat(ws.Path())

	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			stat, err = os.Stat(ws.Path())
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
		log.Fatal(fmt.Sprintf("Workspace requires a dictionary: %v", ws.Path()))
		os.Exit(1)
	}

	_, err = os.Open(ws.Path())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ws.Config = models.GetConfig(ws.path)
	return &ws
}

func (ws *Workspace) GetWd(cleanUp bool, project string) *WorkDir {
	return NewWorkDir(cleanUp, path.Join(ws.Path(), project))
}

func (ws *Workspace) Path() string {
	return ws.path
}
