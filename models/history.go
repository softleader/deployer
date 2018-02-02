package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"encoding/json"
)

type Histories []Deploy

func (h Histories) Len() int           { return len(h) }
func (h Histories) Less(i, j int) bool { return h[i].Project < h[j].Project }
func (h Histories) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func GetHistories(ws string) (histories Histories, err error) {
	f := history(ws)
	_, err = os.OpenFile(f, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	if len(b) <= 0 {
		return []Deploy{}, nil
	}
	err = json.Unmarshal(b, &histories)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (h *Histories) Push(d *Deploy) {
	i := -1
	for idx, e := range *h {
		if e.Project == d.Project {
			i = idx
			break
		}
	}
	if i >= 0 {
		*h = append((*h)[:i], append([]Deploy{*d}, (*h)[i+1:]...)...)
	} else {
		*h = append(*h, *d)
	}
}

func (h *Histories) SaveTo(ws string) (err error) {
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(history(ws), b, os.ModePerm)
}

func (h *Histories) Delete(i int) {
	*h = append((*h)[:i], (*h)[i+1:]...)
}

func history(ws string) string {
	return filepath.Join(ws, "histories.json")
}
