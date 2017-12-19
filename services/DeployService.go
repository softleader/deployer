package services

import (
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/cmd"
	"bytes"
	"io/ioutil"
	"strings"
	"path"
	"fmt"
	"regexp"
	"strconv"
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
		ds.Wd.RemoveAll()
		ds.Wd.MkdirAll()
	}

	var buf bytes.Buffer
	msg := fmt.Sprintf("\nDeploying '%v'...\n", d.Yaml)
	fmt.Print(msg)
	buf.WriteString(msg)

	gpmDir := "repo"
	group, err := ds.gpmInstall(&buf, gpmDir, &d)
	if err != nil {
		return err.Error()
	}

	var yamls []string
	repo := path.Join(ds.Wd.Path, gpmDir)

	if !group {
		yml, err := ds.genYaml(&buf, repo, "docker-compose.yml", &d)
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
			yml, err := ds.genYaml(&buf, groupRepo, fmt.Sprintf("docker-compose-%v.yml", group.Name()), &d)
			if err != nil {
				return err.Error()
			}
			yamls = append(yamls, yml)
		}
	}

	err = ds.delpoyDocker(&buf, yamls, &d)
	if err != nil {
		return err.Error()
	}

	buf.WriteString("\n")

	return buf.String()
}

func (ds *DeployService) gpmInstall(buf *bytes.Buffer, dir string, d *datamodels.Deploy) (bool, error) {
	cmd, out, err := ds.Gpm.Install(dir, d.Yaml)
	if err != nil {
		return false, err
	}
	buf.WriteString(fmt.Sprintf("$ %v\n", cmd))
	buf.WriteString(out)
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func (ds *DeployService) genYaml(buf *bytes.Buffer, dirname string, outYaml string, d *datamodels.Deploy) (string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return "", err
	}
	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, path.Join(dirname, f.Name()))
		}
	}

	yml := path.Join(dirname, outYaml)
	cmd, out, err := ds.GenYaml.Gen(yml, d, strings.Join(dirs, " "))
	if err != nil {
		return "", err
	}

	err = updateDevPort(out, d)
	if err != nil {
		return "", err;
	}

	buf.WriteString(fmt.Sprintf("$ %v\n", cmd))
	buf.WriteString(out)
	return yml, nil
}

func updateDevPort(out string, d *datamodels.Deploy) error {
	if d.Dev.Addr != "" {
		re, err := regexp.Compile(`Auto publish port from \[\d*\] to \[(\d*)\]`)
		if err != nil {
			return err
		}
		res := re.FindStringSubmatch(out)
		if len(res) > 1 {
			d.Dev.Port, err = strconv.Atoi(res[1])
			if err != nil {
				return err
			}
			d.Dev.Port++
		}
	}
	return nil
}

func (ds *DeployService) delpoyDocker(buf *bytes.Buffer, yamls []string, d *datamodels.Deploy) error {
	for _, yaml := range yamls {
		stack := d.Project
		if d.Dev.Addr != "" {
			stack += "_" + strconv.Itoa(d.Dev.Port)
		}
		cmd, out, err := ds.DockerStack.Deploy(stack, yaml)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf("$ %v\n", cmd))
		buf.WriteString(out)
	}
	return nil
}

func (ds *DeployService) Delete(stack string) string {
	_, out, err := ds.DockerStack.RmLike(stack)
	if err != nil {
		return err.Error()
	} else {
		return out
	}
}
