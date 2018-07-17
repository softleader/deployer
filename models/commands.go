package models

import (
	"strings"
	"strconv"
)

type DockerStatsNoStream struct {
	Name      string
	TruncName string
	CPUPerc   string
	MemUsage  string
	MemPerc   string
	NetIO     string
	BlockIO   string
}

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

type DockerServiceLs struct {
	Replicas string
}

type StackName struct {
	Project         string
	Port            int
	Group           string
	IsDevEnabled    bool
	IsPublishedPort bool
}

func NewStackName(in string) (n StackName) {
	s := strings.Split(in, "-")
	n.Project = s[0]
	if len(s) > 1 {
		n.IsDevEnabled = true
		p, err := strconv.Atoi(s[1])
		if err == nil {
			n.Port = p
			n.IsPublishedPort = true
		}
	}
	if len(s) > 2 {
		n.Group = s[2]
	}
	return
}

type Replicas struct {
	Up    int
	Down  int
	Total int
}

func NewReplicas(in string) (r Replicas) {
	s := strings.Split(in, "/")
	r.Up, _ = strconv.Atoi(s[0])
	r.Total, _ = strconv.Atoi(s[1])
	r.Down = r.Total - r.Up
	return
}

func (r Replicas) LooksGood() bool {
	return r.Up == r.Total
}

// convert Replicas to integer
func (d *DockerServiceLs) Rtoi() (r Replicas) {
	r = NewReplicas(d.Replicas)
	return
}

// convert Replicas to integer
func (d *DockerStackServices) Rtoi() (r Replicas) {
	r = NewReplicas(d.Replicas)
	return
}

func NewDockerNodeLs(in string) (m DockerNodeLs) {
	s := strings.Split(in, ";")
	m.Hostname = s[0]
	m.Status = s[1]
	m.Availability = s[2]
	return
}

func NewDockerStackServices(in string) (m DockerStackServices) {
	s := strings.Split(in, ";")
	m.Id = s[0]
	m.Name = s[1]
	m.Mode = s[2]
	m.Replicas = s[3]
	m.Image = s[4]
	m.Ports = s[5]
	return
}

func NewDockerStackLs(in string) (m DockerStackLs) {
	s := strings.Split(in, ";")
	m.Name = s[0]
	m.Services = s[1]
	return
}

func NewDockerStatsNoSteam(in string) (m DockerStatsNoStream) {
	s := strings.Split(in, ";")
	m.Name = s[0]
	name := strings.Split(m.Name, ".")
	m.TruncName = strings.Join(name[:len(name)-1 ], ".")
	m.CPUPerc = s[1]
	m.MemUsage = s[2]
	m.MemPerc = s[3]
	m.NetIO = s[4]
	m.BlockIO = s[5]
	return
}

func NewDockerServicePs(in string) (m DockerServicePs) {
	s := strings.Split(in, ";")
	m.Id = s[0]
	m.Name = s[1]
	m.Image = strings.Split(s[2], "@sha256")[0]
	m.Node = s[3]
	m.DesiredState = s[4]
	m.CurrentState = s[5]
	m.Error = s[6]
	return
}

func NewDockerServiceLs(in string) (m DockerServiceLs) {
	s := strings.Split(in, ";")
	m.Replicas = s[0]
	return
}
