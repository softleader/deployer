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
	cmd.Workspace
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
	wd := ds.Workspace.GetWd(d.CleanUp, d.Project)
	opts := cmd.Options{Ctx: ctx, Pwd: wd.Path}
	d.Dev.PublishPort = d.Dev.Port

	(*ctx).StreamWriter(pipe.Printf("\nDeploying '%v'...\n", d.Yaml))

	gpmDir := "repo"
	group, err := ds.gpmInstall(&opts, gpmDir, &d)
	if err != nil {
		return err
	}
	var yamls []datamodels.Yaml
	repo := path.Join(wd.Path, gpmDir)

	if !group {
		yml := path.Join(repo, "docker-compose.yml")
		dirs, err := collectDirs(repo)
		if err != nil {
			return err
		}
		err = ds.genYaml(&opts, dirs, yml, &d)
		if err != nil {
			return err
		}
		yamls = append(yamls, datamodels.Yaml{
			Group: "",
			Path:  yml,
		})
	} else {
		// 限定一層的 group
		groups, err := ioutil.ReadDir(repo)
		if err != nil {
			return err
		}

		deployGroups := map[string][]string{}
		for _, group := range groups {
			if d.Group != "" && !d.GroupContains(group.Name()) {
				(*ctx).StreamWriter(pipe.Printf("Skip deploying group [%v] because it does not match any of '%v'\n", group.Name(), d.Group))
				continue
			}
			dirs, err := collectDirs(path.Join(repo, group.Name()))
			if err != nil {
				return err
			}
			deployGroups[group.Name()] = dirs
		}

		if d.FlatGroup {
			var flat []string
			for _, dirs := range deployGroups {
				for _, d := range dirs {
					flat = append(flat, d)
				}
			}
			yml := path.Join(repo, "docker-compose.yml")
			err := ds.genYaml(&opts, flat, yml, &d)
			if err != nil {
				return err
			}
			yamls = append(yamls, datamodels.Yaml{
				Group: "",
				Path:  yml,
			})
		} else {
			for group, dirs := range deployGroups {
				yml := path.Join(repo, group, fmt.Sprintf("docker-compose-%v.yml", group))
				err := ds.genYaml(&opts, dirs, yml, &d)
				if err != nil {
					return err
				}
				yamls = append(yamls, datamodels.Yaml{
					Group: group,
					Path:  yml,
				})
			}
		}
	}

	err = wd.CopyToDeployedDir(yamls)
	if err != nil {
		return err
	}

	err = wd.CompressDeployedDir()
	if err != nil {
		return err
	}

	err = ds.deployDocker(&opts, yamls, &d)
	if err != nil {
		return err
	}

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

func collectDirs(p string) ([]string, error) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return []string{}, err
	}
	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, path.Join(p, f.Name()))
		}
	}
	return dirs, nil
}

func (ds *DeployService) genYaml(opts *cmd.Options, dirs []string, output string, d *datamodels.Deploy) error {
	_, out, err := ds.GenYaml.Gen(opts, output, d, strings.Join(dirs, " "))
	if err != nil {
		return err
	}
	err = updateDevPort(out, d)
	if err != nil {
		return err
	}
	(*opts.Ctx).StreamWriter(pipe.Print(out))
	return nil
}

func updateDevPort(out string, d *datamodels.Deploy) error {
	if d.Dev.IpAddress != "" {
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

func (ds *DeployService) deployDocker(opts *cmd.Options, yamls []datamodels.Yaml, d *datamodels.Deploy) error {
	for _, y := range yamls {
		stack := []string{d.Project}
		if d.Dev.IpAddress != "" {
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
	_, out, err := ds.DockerStack.RmStackLike(stack)
	return out, err
}

func (ds *DeployService) DeleteService(service string) (string, error) {
	_, out, err := ds.DockerStack.RmService(service)
	return out, err
}
