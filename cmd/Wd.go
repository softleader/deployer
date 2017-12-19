package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
)

type Wd struct {
	Path string
}

func NewWd(dir string) Wd {
	if dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/wd")
	}
	wd := Wd{Path: dir}
	fmt.Printf("Setting up working directory to '%v'\n", dir)

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
		log.Fatal(fmt.Sprintf("requires a dictionary: %v", wd.Path))
		os.Exit(1)
	}

	_, err = os.Open(wd.Path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return wd
}

func (wd Wd) checkWd() (bool, error) {
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

func (wd Wd) RemoveAll() Wd {
	os.RemoveAll(wd.Path)
	return wd
}

func (wd Wd) MkdirAll() Wd {
	os.MkdirAll(wd.Path, os.ModeDir|os.ModePerm)
	return wd
}
