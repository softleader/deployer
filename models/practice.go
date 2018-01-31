package models

import (
	"io/ioutil"
	"os"
)

import (
	"path/filepath"
)

func ReadPractices(ws string) (content string, err error) {
	f := practices(ws)
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
	return ioutil.WriteFile(practices(ws), []byte(content), os.ModePerm)
}

func practices(ws string) string {
	return filepath.Join(ws, "best-practices.md")
}
