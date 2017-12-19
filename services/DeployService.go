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

	var buffer bytes.Buffer
	msg := fmt.Sprintf("\nDeploying '%v'...\n", d.Yaml)
	fmt.Print(msg)
	buffer.WriteString(msg)

	gpmDir := "repo"
	group, err := ds.gpmInstall(&buffer, gpmDir, &d)
	if err != nil {
		return err.Error()
	}

	var yamls []string
	repo := path.Join(ds.Wd.Path, gpmDir)

	if !group {
		yml, err := ds.genYaml(&buffer, repo, "docker-compose.yml", &d)
		if err != nil {
			return err.Error()
		}
		yamls = append(yamls, yml)
	} else {
		// 目前只支援一層的 group..
		groups, err := ioutil.ReadDir(repo)
		if err != nil {
			return err.Error()
		}
		for _, group := range groups {
			groupRepo := path.Join(repo, group.Name())
			yml, err := ds.genYaml(&buffer, groupRepo, fmt.Sprintf("docker-compose-%v.yml", group.Name()), &d)
			if err != nil {
				return err.Error()
			}
			yamls = append(yamls, yml)
		}
	}

	err = ds.delpoyDocker(&buffer, yamls, &d)
	if err != nil {
		return err.Error()
	}

	buffer.WriteString("\n")

	return buffer.String()
}

func (ds *DeployService) gpmInstall(buffer *bytes.Buffer, dir string, d *datamodels.Deploy) (bool, error) {
	cmd, out, err := ds.Gpm.Install(dir, d.Yaml)
	if err != nil {
		return false, err
	}
	buffer.WriteString(fmt.Sprintf("$ %v", cmd))
	buffer.WriteString(out)
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func (ds *DeployService) genYaml(buffer *bytes.Buffer, dirname string, outYaml string, d *datamodels.Deploy) (string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return "", err
	}
	var dirs []string
	for _, f := range files {
		dirs = append(dirs, path.Join(dirname, f.Name()))
	}

	yml := path.Join(dirname, outYaml)
	cmd, out, err := ds.GenYaml.Gen(yml, *d, strings.Join(dirs, " "))
	if err != nil {
		return "", err
	}

	d.Dev.Port = retrieveDevPort(out)

	buffer.WriteString(fmt.Sprintf("$ %v", cmd))
	buffer.WriteString(out)
	return yml, nil
}

func retrieveDevPort(out string) int {
	return 1
}

func (ds *DeployService) delpoyDocker(buffer *bytes.Buffer, yamls []string, d *datamodels.Deploy) error {
	for _, yaml := range yamls {
		cmd, out, err := ds.DockerStack.Deploy(d.Project, yaml)
		if err != nil {
			return err
		}
		buffer.WriteString(fmt.Sprintf("$ %v", cmd))
		buffer.WriteString(out)
	}
	return nil
}

func (ds *DeployService) Delete(stack string) string {
	_, out, err := ds.DockerStack.Rm(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}
