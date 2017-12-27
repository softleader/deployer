package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
)

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
			wd.MkdirAll("")
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

func (ws *Ws) Pwd(project string) string {
	return path.Join(ws.path, project)
}

func (wd *Ws) RemoveAll(project string) {
	os.RemoveAll(wd.Pwd(project))
}

func (wd *Ws) MkdirAll(project string) {
	os.MkdirAll(wd.Pwd(project), os.ModeDir|os.ModePerm)
}
