package services

import (
	"github.com/softleader/deployer/datamodels"
	"github.com/softleader/deployer/cmd"
	"io/ioutil"
	"strings"
	"path"
	"fmt"
	"regexp"
	"strconv"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/pipe"
)

type DeployService struct {
	cmd.DockerStack
	cmd.Gpm
	cmd.GenYaml
	cmd.Ws
}

type compose struct {
	group string
	yaml  string
}

func (ds *DeployService) GetAll() ([][]string, error) {
	_, out, err := ds.DockerStack.Ls()
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, strings.Split(line, ";"))
		}
	}
	return s, err
}

func (ds *DeployService) GetServices(stack string) ([][]string, error) {
	_, out, err := ds.DockerStack.Services(stack)
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			s = append(s, strings.Split(line, ";"))
		}
	}
	return s, err
}

func (ds *DeployService) Ps(id string) ([][]string, error) {
	_, out, err := ds.DockerStack.Ps(id)
	lines := strings.Split(out, "\n")
	var s [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fields := strings.Split(line, ";")
			fields[2] = strings.Split(fields[2], "@sha256")[0]
			s = append(s, fields)
		}
	}
	return s, err
}

func (ds *DeployService) Deploy(ctx *iris.Context, d datamodels.Deploy) error {
	if d.CleanUp {
		ds.Ws.RemoveAll()
		ds.Ws.MkdirAll()
	}

	d.Dev.PublishPort = d.Dev.Port

	(*ctx).StreamWriter(pipe.Printf("\nDeploying '%v'...\n", d.Yaml))

	gpmDir := "repo"
	group, err := ds.gpmInstall(ctx, gpmDir, &d)
	if err != nil {
		return err
	}
	var c []compose
	repo := path.Join(ds.Ws.Path, gpmDir)

	if !group {
		yml, err := ds.genYaml(ctx, repo, "docker-compose.yml", &d)
		if err != nil {
			return err
		}
		c = append(c, compose{
			group: "",
			yaml:  yml,
		})
	} else {
		// 目前只支援一層的 group..
		groups, err := ioutil.ReadDir(repo)
		if err != nil {
			return err
		}
		for _, group := range groups {
			groupRepo := path.Join(repo, group.Name())
			yml, err := ds.genYaml(ctx, groupRepo, fmt.Sprintf("docker-compose-%v.yml", group.Name()), &d)
			if err != nil {
				return err
			}
			c = append(c, compose{
				group: group.Name(),
				yaml:  yml,
			})
		}
	}

	err = ds.deployDocker(ctx, c, &d)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DeployService) gpmInstall(ctx *iris.Context, dir string, d *datamodels.Deploy) (bool, error) {
	_, out, err := ds.Gpm.Install(ctx, dir, d.Yaml)
	if err != nil {
		return false, err
	}
	(*ctx).StreamWriter(pipe.Print(out))
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func (ds *DeployService) genYaml(ctx *iris.Context, dirname string, outYaml string, d *datamodels.Deploy) (string, error) {
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
	_, out, err := ds.GenYaml.Gen(ctx, yml, d, strings.Join(dirs, " "))
	if err != nil {
		return "", err
	}

	err = updateDevPort(out, d)
	if err != nil {
		return "", err;
	}

	(*ctx).StreamWriter(pipe.Print(out))
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
			d.Dev.PublishPort, err = strconv.Atoi(res[1])
			if err != nil {
				return err
			}
			d.Dev.PublishPort++
		}
	}
	return nil
}

func (ds *DeployService) deployDocker(ctx *iris.Context, composes []compose, d *datamodels.Deploy) error {
	for _, c := range composes {
		stack := []string{d.Project}
		if d.Dev.Addr != "" {
			stack = append(stack, strconv.Itoa(d.Dev.Port))
		}
		if c.group != "" {
			stack = append(stack, c.group)
		}
		ds.DockerStack.Deploy(ctx, strings.Join(stack, "-"), c.yaml)

	}
	return nil
}

func (ds *DeployService) Delete(stack string) (string, error) {
	_, out, err := ds.DockerStack.RmLike(stack)
	return out, err
}
