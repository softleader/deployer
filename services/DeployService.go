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
		ds.Ws.RemoveAll(d.Project)
		ds.Ws.MkdirAll(d.Project)
	}

	opts := cmd.Options{Ctx: ctx, Pwd: ds.Pwd(d.Project)}
	d.Dev.PublishPort = d.Dev.Port

	(*ctx).StreamWriter(pipe.Printf("\nDeploying '%v'...\n", d.Yaml))

	gpmDir := "repo"
	group, err := ds.gpmInstall(&opts, gpmDir, &d)
	if err != nil {
		return err
	}
	y := datamodels.Yamls{}
	repo := path.Join(ds.Ws.Pwd(d.Project), gpmDir)

	if !group {
		yml, err := ds.genYaml(&opts, repo, "docker-compose.yml", &d)
		if err != nil {
			return err
		}
		y = append(y, datamodels.Yaml{
			Group: "",
			Path:  yml,
		})
	} else {
		// 目前只支援一層的 group..
		groups, err := ioutil.ReadDir(repo)
		if err != nil {
			return err
		}
		for _, group := range groups {
			if d.Group != "" && !d.GroupContains(group.Name()) {
				(*ctx).StreamWriter(pipe.Printf("Skip deploying group [%v] because it does not match any of '%v'\n", group.Name(), d.Group))
				continue
			}
			groupRepo := path.Join(repo, group.Name())
			yml, err := ds.genYaml(&opts, groupRepo, fmt.Sprintf("docker-compose-%v.yml", group.Name()), &d)
			if err != nil {
				return err
			}
			y = append(y, datamodels.Yaml{
				Group: group.Name(),
				Path:  yml,
			})
		}
	}

	err = ds.deployDocker(&opts, y, &d)
	if err != nil {
		return err
	}

	y.ZipTo(opts.Pwd)

	return nil
}

func (ds *DeployService) gpmInstall(opts *cmd.Options, dir string, d *datamodels.Deploy) (bool, error) {
	_, out, err := ds.Gpm.Install(opts, dir, d.Yaml)
	if err != nil {
		return false, err
	}
	(*opts.Ctx).StreamWriter(pipe.Print(out))
	return strings.Contains(out, "Detected groups in YAML dependencies!"), nil
}

func (ds *DeployService) genYaml(opts *cmd.Options, dirname string, outYaml string, d *datamodels.Deploy) (string, error) {
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
	_, out, err := ds.GenYaml.Gen(opts, yml, d, strings.Join(dirs, " "))
	if err != nil {
		return "", err
	}

	err = updateDevPort(out, d)
	if err != nil {
		return "", err
	}

	(*opts.Ctx).StreamWriter(pipe.Print(out))
	return yml, nil
}

func updateDevPort(out string, d *datamodels.Deploy) error {
	if d.Dev.Hostname != "" {
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

func (ds *DeployService) deployDocker(opts *cmd.Options, yamls datamodels.Yamls, d *datamodels.Deploy) error {
	for _, y := range yamls {
		stack := []string{d.Project}
		if d.Dev.Hostname != "" {
			stack = append(stack, strconv.Itoa(d.Dev.Port))
		}
		if y.Group != "" {
			stack = append(stack, y.Group)
		}
		ds.DockerStack.Deploy(opts, strings.Join(stack, "-"), y.Path)

	}
	return nil
}

func (ds *DeployService) DeleteStack(stack string) (string, error) {
	_, out, err := ds.DockerStack.RmStack(stack)
	return out, err
}

func (ds *DeployService) DeleteService(service string) (string, error) {
	_, out, err := ds.DockerStack.RmService(service)
	return out, err
}
