package models

import (
	"strings"
	"fmt"
)

type DockerNodeLs struct {
	Hostname     string
	Status       string
	Availability string
}

type DockerStackServices struct {
	Id       string
	Name     string
	Mode     string
	Replicas string
	Image    string
	Ports    string
}

type DockerStackLs struct {
	Name     string
	Services string
	Uptime   string
}

type DockerServicePs struct {
	Id           string
	Name         string
	Image        string
	Node         string
	DesiredState string
	CurrentState string
	Error        string
}

func NewDockerNodeLs(out string) (m DockerNodeLs) {
	s := strings.Split(out, ";")
	fmt.Println(s)
	m.Hostname = s[0]
	m.Status = s[1]
	m.Availability = s[2]
	return
}

func NewDockerStackServices(out string) (m DockerStackServices) {
	s := strings.Split(out, ";")
	m.Id = s[0]
	m.Name = s[1]
	m.Mode = s[2]
	m.Replicas = s[3]
	m.Image = s[4]
	m.Ports = s[5]
	return
}

func NewDockerStackLs(out string) (m DockerStackLs) {
	s := strings.Split(out, ";")
	m.Name = s[0]
	m.Services = s[1]
	return
}

func NewDockerServicePs(out string) (m DockerServicePs) {
	s := strings.Split(out, ";")
	m.Id = s[0]
	m.Name = s[1]
	m.Image = strings.Split(s[2], "@sha256")[0]
	m.Node = s[3]
	m.DesiredState = s[4]
	m.CurrentState = s[5]
	m.Error = s[6]
	return
}
