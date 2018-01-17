package models

import (
	"io/ioutil"
	"os"
)

import (
	"path/filepath"
)

var (
	filename = "best-practices.md"
)

func ReadPractices(ws string) (content string, err error) {
	f := file(ws)
	_, err = os.OpenFile(f, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SavePractices(ws string, content string) (err error) {
	return ioutil.WriteFile(file(ws), []byte(content), os.ModePerm)
}

func file(ws string) string {
	return filepath.Join(ws, filename)
}
