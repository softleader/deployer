package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"encoding/json"
)

func GetHistories(ws string) (histories []Deploy, err error) {
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

func SaveHistory(ws string, d *Deploy) (err error) {
	histories, err := GetHistories(ws)
	if err != nil {
		return err
	}

	i := 0
	for idx, h := range histories {
		if h.Project == d.Project {
			i = idx
			break
		}
	}
	if i > 0 {
		histories = append(histories[:i], append([]Deploy{*d}, histories[i+1:]...)...)
	} else {
		histories = append(histories, *d)
	}
	return writeHistories(ws, histories)
}

func RemoveHistory(ws string, i int) (err error) {
	histories, err := GetHistories(ws)
	if err != nil {
		return err
	}
	histories = append(histories[:i], histories[i+1:]...)
	return writeHistories(ws, histories)
}

func writeHistories(ws string, h []Deploy) (err error) {
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(history(ws), b, os.ModePerm)
}

func history(ws string) string {
	return filepath.Join(ws, "histories.json")
}
