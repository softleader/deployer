package datamodels

import "strings"

type Deploy struct {
	CleanUp   bool   `json:cleanUp`
	Project   string `json:"project"`
	Dev       Dev    `json:"dev"`
	Volume0   string `json:"volume0"`
	Net0      string `json:"net0"`
	Yaml      string `json:"yaml"`
	Group     string `json:"group"`
	FlatGroup bool   `json:"flatGroup"`
	Silently  bool   `json:"silently"`
}

func (d *Deploy) GroupContains(group string) bool {
	groups := strings.Split(d.Group, ",")
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}
