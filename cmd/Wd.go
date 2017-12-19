package cmd

import (
	"os"
	"path"
	"log"
	"fmt"
	"os/user"
)

type Wd struct {
	Path string
}

func NewWd(dir string) Wd {
	if dir == "" {
		dir, _ := os.Getwd()
		dir = path.Join(dir, "/wd")
	}
	wd := Wd{Path: dir}

	stat, err := os.Stat(wd.Path)

	if err != nil {
		if os.IsNotExist(err) {
			wd.MkdirAll()
		} else {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	if !stat.IsDir() {
		log.Fatal(fmt.Sprintf("'%v' requires a dir", dir))
		os.Exit(1)
	}

	if stat.Mode().Perm() != os.ModePerm {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.Fatal(fmt.Sprintf("'%v' requires permission '%v' to '%v'", dir, os.ModePerm, u.Name))
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
