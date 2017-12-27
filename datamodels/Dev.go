package datamodels

import (
	"strconv"
	"strings"
)

type Dev struct {
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"` // 最初傳進來的 port
	Ignore      string `json:"ignore"`
	PublishPort int    `json:"-"` // 紀錄當前 publish 最後的 port
}

func (d *Dev) String() string {
	s := []string{"--dev-hostname " + d.Hostname}
	if d.PublishPort > 0 {
		s = append(s, "--dev-port "+strconv.Itoa(d.PublishPort))
	}
	if d.Ignore != "" {
		s = append(s, "--dev-ignore "+d.Ignore)
	}
	return strings.Join(s, " ")
}
