package models

import (
	"strings"
	"path/filepath"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Deploy Deploy            `json:"deploy"`
	Navbar map[string]string `json:"navbar"`
}

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
	b, _ := yaml.Marshal(&d)
	return string(b)
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

func GetConfig(ws string) Config {
	config := filepath.Join(ws, "config.yaml")
	if _, err := os.Stat(config); os.IsNotExist(err) {
		dft := Deploy{
			CleanUp: true,
			Style:   "swarm",
			Dev: Dev{
				IpAddress: "192.168.1.60",
				Port:      0,
				Ignore:    "elasticsearch,kibana,logstash,redis,eureka,softleader-config-server,ldap-server",
			},
			Yaml:    "github:softleader/softleader-package/",
			Extend:  "",
			Volume0: "",
			Net0:    "",
			Group:   "",
		}
		cfg := Config{
			Deploy: dft,
			Navbar: make(map[string]string),
		}
		cfg.Navbar["REST API"] = "https://github.com/softleader/deployer#rest-api"
		b, _ := yaml.Marshal(cfg)
		ioutil.WriteFile(config, b, os.ModePerm)
		return cfg
	} else {
		b, _ := ioutil.ReadFile(config)
		cfg := Config{}
		yaml.Unmarshal(b, &cfg)
		return cfg
	}

}
