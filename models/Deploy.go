package models

import "strings"

type Deploy struct {
	CleanUp   bool   `json:cleanUp`
	Project   string `json:"project"`
	Style     string `json:"style"`
	Dev       Dev    `json:"dev"`
	Volume0   string `json:"volume0"`
	Net0      string `json:"net0"`
	Yaml      string `json:"yaml"`
	Extend    string `json:"extend"`
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
