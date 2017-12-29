package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
)

// workspace, 一個 golang app 會有唯一的一個 workspace
type Ws struct {
	path string
}

func NewWs(dir string) *Ws {
	if dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/workspace")
	}
	wd := Ws{path: dir}
	fmt.Printf("Setting up workspace to '%v'\n", dir)

	stat, err := os.Stat(wd.path)

	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			stat, err = os.Stat(wd.path)
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
		log.Fatal(fmt.Sprintf("Workspace requires a dictionary: %v", wd.path))
		os.Exit(1)
	}

	_, err = os.Open(wd.path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &wd
}

func (wd *Ws) checkWs() (bool, error) {
	stat, err := os.Stat(wd.path)
	if err == nil {
		if stat.IsDir() && stat.Mode().Perm() == os.ModePerm {
			return true, nil
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (ws *Ws) GetWd(project string) Wd {
	return Wd{Path: path.Join(ws.path, project)}
}
