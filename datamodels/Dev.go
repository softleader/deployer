package datamodels

import "strconv"

type Dev struct {
	Addr        string `json:"addr"`
	Port        int    `json:"port"` // 最初傳進來的 port
	PublishPort int    `json:"-"`    // 紀錄當前 publish 最後的 port
}

func (d *Dev) String() string {
	if d.PublishPort <= 0 {
		return d.Addr
	}
	return d.Addr + "/" + strconv.Itoa(d.PublishPort)
}
