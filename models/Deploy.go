package models

import (
	"strings"
	"encoding/json"
)

type Deploy struct {
	Project   string `json:"project"`
	Yaml      string `json:"yaml"`
	Extend    string `json:"extend,omitempty"`
	Dev       Dev    `json:"dev,omitempty"`
	Net0      string `json:"net0,omitempty"`
	Volume0   string `json:"volume0,omitempty"`
	Group     string `json:"group,omitempty"`
	FlatGroup bool   `json:"flatGroup"`
	Style     string `json:"style"`
	CleanUp   bool   `json:"cleanUp"`
	Silently  bool   `json:"silently"`
}

func (d *Deploy) String() string {
	b, _ := json.MarshalIndent(&d, "", "  ")
	return string(b);
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

func NewDefaultDeploy() *Deploy {
	return &Deploy{
		CleanUp: true,
		Dev: Dev{
			IpAddress: "192.168.1.60",
			Port:      0,
			Ignore:    "elasticsearch,kibana,logstash,redis,eureka,softleader-config-server",
		},
		Yaml:    "github:softleader/softleader-package/",
		Extend:  "",
		Volume0: "",
		Net0:    "",
		Group:   "",
	}
}
