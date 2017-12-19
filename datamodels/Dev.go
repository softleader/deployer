package datamodels

import "strconv"

type Dev struct {
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
	PublishPort int
}

func (d *Dev) String() string {
	if d.Port <= 0 {
		return d.Addr
	}
	return d.Addr + "/" + strconv.Itoa(d.Port)
}
