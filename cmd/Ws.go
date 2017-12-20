package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
)

type Ws struct {
	Path string
}

func NewWs(dir string) *Ws {
	if dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/workspace")
	}
	wd := Ws{Path: dir}
	fmt.Printf("Setting up workspace to '%v'\n", dir)

	stat, err := os.Stat(wd.Path)

	if err != nil {
		if os.IsNotExist(err) {
			wd.MkdirAll()
			stat, err = os.Stat(wd.Path)
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
		log.Fatal(fmt.Sprintf("Workspace requires a dictionary: %v", wd.Path))
		os.Exit(1)
	}

	_, err = os.Open(wd.Path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &wd
}

func (wd *Ws) checkWs() (bool, error) {
	stat, err := os.Stat(wd.Path)
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

func (wd *Ws) RemoveAll() {
	os.RemoveAll(wd.Path)
}

func (wd *Ws) MkdirAll() {
	os.MkdirAll(wd.Path, os.ModeDir|os.ModePerm)
}
