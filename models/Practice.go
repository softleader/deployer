package models

import (
	"io/ioutil"
	"os"
)

import (
	"encoding/json"
	"path/filepath"
)

var (
	filename = "practices.json"
)

type Practice []string

func ReadFromFile(ws string) (p Practice, err error) {
	f := file(ws)
	_, err = os.OpenFile(f, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	p = Practice{}
	if len(b) > 0 {
		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Practice) SaveToFile(ws string) (err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file(ws), b, os.ModePerm)
}

func file(ws string) string {
	return filepath.Join(ws, filename)
}
