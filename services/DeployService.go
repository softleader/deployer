package services

import (
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/cmd"
	"bytes"
	"io/ioutil"
	"strings"
	"path"
	"fmt"
)

type DeployService struct {
	cmd.DockerStack
	cmd.Gpm
	cmd.GenYaml
	cmd.Wd
}

func (ds *DeployService) GetAll() string {
	_, out, err := ds.DockerStack.Ls()
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (ds *DeployService) GetServices(stack string) string {
	_, out, err := ds.DockerStack.Services(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}

func (ds *DeployService) Deploy(d datamodels.Deploy) string {
	if d.CleanUp {
		ds.Wd.RemoveAll().MkdirAll()
	}

	var resp bytes.Buffer
	for _, y := range d.Yaml {
		resp.WriteString(ds.deploy(y, d))
	}

	return resp.String()
}

func (ds *DeployService) deploy(y string, d datamodels.Deploy) string {
	var resp bytes.Buffer
	msg := fmt.Sprintf("\nDeploying '%v'...\n", y)
	fmt.Print(msg)
	resp.WriteString(msg)

	// gpm install
	installDir := "repo"
	cmd, out, err := ds.Gpm.Install(installDir, y)
	if err != nil {
		return err.Error()
	}
	resp.WriteString(fmt.Sprintf("$ %v", cmd))
	resp.WriteString(out)

	// gen-yaml

	var yamls []string
	repo := path.Join(ds.Wd.Path, installDir)

	if !strings.Contains(out, "Detected groups in YAML dependencies!") {
		yml, cmd, out, err := ds.genYaml(repo, "docker-compose.yml", d)
		if err != nil {
			return err.Error()
		}
		resp.WriteString(fmt.Sprintf("$ %v", cmd))
		resp.WriteString(out)
		yamls = append(yamls, yml)
	} else {

		// 目前只支援一層的 group..

		groups, err := ioutil.ReadDir(repo)
		if err != nil {
			return err.Error()
		}

		for _, group := range groups {
			groupRepo := path.Join(repo, group.Name())
			yml, cmd, out, err := ds.genYaml(groupRepo, fmt.Sprintf("docker-compose-%v.yml", group.Name()), d)
			if err != nil {
				return err.Error()
			}
			resp.WriteString(fmt.Sprintf("$ %v", cmd))
			resp.WriteString(out)
			yamls = append(yamls, yml)
		}
	}

	// docker stack deploy

	for _, yaml := range yamls {
		cmd, out, err = ds.DockerStack.Deploy(d.Project, yaml)
		if err != nil {
			return err.Error()
		}
		resp.WriteString(fmt.Sprintf("$ %v", cmd))
		resp.WriteString(out)
	}

	resp.WriteString("\n")

	return resp.String()
}

func (ds *DeployService) genYaml(dirname string, yaml string, d datamodels.Deploy) (string, string, string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return "", "", "", err
	}
	var dirs []string
	for _, f := range files {
		dirs = append(dirs, path.Join(dirname, f.Name()))
	}

	yml := path.Join(dirname, "docker-compose.yml")
	cmd, out, err := ds.GenYaml.Gen(yml, d, strings.Join(dirs, " "))
	if err != nil {
		return "", "", "", err
	}

	return yml, cmd, out, nil
}

func (ds *DeployService) Delete(stack string) string {
	_, out, err := ds.DockerStack.Rm(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}
