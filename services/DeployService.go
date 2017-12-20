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
	"io"
	"github.com/kataras/iris"
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

func (ds *DeployService) GetAll() (string, error) {
	_, out, err := ds.DockerStack.Ls()
	return out, err
}

func (ds *DeployService) GetServices(stack string) (string, error) {
	_, out, err := ds.DockerStack.Services(stack)
	return out, err
}

func (ds *DeployService) Ps(id string) (string, error) {
	_, out, err := ds.DockerStack.Ps(id)
	return out, err
}

func (ds *DeployService) Deploy(ctx *iris.Context, d datamodels.Deploy) error {
	if d.CleanUp {
		ds.Ws.RemoveAll()
		ds.Ws.MkdirAll()
	}

	d.Dev.PublishPort = d.Dev.Port

	msg := fmt.Sprintf("\nDeploying '%v'...\n", d.Yaml)
	fmt.Print(msg)
	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, msg)
		return false
	})

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

	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprintln(w)
		return false
	})

	return nil
}

func (ds *DeployService) gpmInstall(ctx *iris.Context, dir string, d *datamodels.Deploy) (bool, error) {
	cmd, out, err := ds.Gpm.Install(dir, d.Yaml)
	if err != nil {
		return false, err
	}
	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, "$ %v\n", cmd)
		return false
	})
	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprint(w, out)
		return false
	})
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
	cmd, out, err := ds.GenYaml.Gen(yml, d, strings.Join(dirs, " "))
	if err != nil {
		return "", err
	}

	err = updateDevPort(out, d)
	if err != nil {
		return "", err;
	}

	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, "$ %v\n", cmd)
		return false
	})
	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprint(w, out)
		return false
	})
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
