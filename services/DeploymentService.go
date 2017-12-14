package services

import (
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/cmd"
	"bytes"
	"io/ioutil"
	"strings"
	"fmt"
	"path"
)

type DeploymentService struct {
	cmd.DockerStack
	cmd.Gpm
	cmd.GenYaml
	cmd.Wd
}

func NewDeploymentService() DeploymentService {
	return DeploymentService{
		DockerStack: cmd.NewDockerStack(),
		Gpm:         cmd.NewGpm(),
		GenYaml:     cmd.NewGenYaml(),
		Wd:          cmd.NewWd(),
	}
}

func (ds *DeploymentService) GetAll() string {
	out, err := ds.DockerStack.Ls()
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (ds *DeploymentService) GetServices(stack string) string {
	out, err := ds.DockerStack.Services(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (ds *DeploymentService) Deploy(d datamodels.Deployment) string {
	if strings.Contains(d.Yaml, "github:") {
		ds.Wd.RemoveAll().MkdirAll()
	}

	var resp bytes.Buffer

	// gpm install
	installDir := "repo"
	out, err := ds.Gpm.Install(installDir, d.Yaml)
	if err != nil {
		return err.Error()
	}
	resp.WriteString(out)

	// gen-yaml

	var yamls []string
	repo := path.Join(ds.Wd.Path, installDir)

	if !strings.Contains(out, "Detected groups in YAML dependencies!") {
		files, err := ioutil.ReadDir(repo)
		if err != nil {
			return err.Error()
		}
		var dirs []string
		for _, f := range files {
			dirs = append(dirs, path.Join(repo, f.Name()))
		}

		yml := path.Join(repo, "docker-compose.yml")
		out, err := ds.GenYaml.Gen(yml, d, strings.Join(dirs, " "))
		if err != nil {
			return err.Error()
		}
		resp.WriteString(out)
		yamls = append(yamls, yml)
	} else {

	}

	// docker stack deploy

	for _, yaml := range yamls {
		//out, err = ds.DockerStack.Deploy(d.Project, yaml)
		//if err != nil {
		//	return err.Error()
		//}
		//resp.WriteString(out)
		fmt.Println("deploy " + yaml)
	}

	return resp.String()
}

func (ds *DeploymentService) Delete(stack string) string {
	out, err := ds.DockerStack.Rm(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}
