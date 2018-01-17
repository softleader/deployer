package routes

import (
	"github.com/kataras/iris"
	"strings"
	"time"
	"github.com/softleader/deployer/models"
	"fmt"
	"strconv"
	"github.com/softleader/deployer/pipe"
	"encoding/json"
	"github.com/softleader/deployer/cmd"
	"path"
	"io/ioutil"
	"github.com/softleader/deployer/app"
)

type StackRoutes struct {
	app.Workspace
	cmd.DockerStack
	cmd.DockerService
	cmd.Gpm
	cmd.GenYaml
}

func (r *StackRoutes) ListStack(ctx iris.Context) {
	out, err := r.DockerStack.Ls()
	if err != nil {
		out = append(out, []string{err.Error()})
	}
	stacks := make(map[string][][]string)
	for _, line := range out {
		splited := strings.Split(line[0], "-")
		key := splited[0]
		if len(splited) > 1 {
			if publishedPort(splited[1]) { // 有 publish port 可視為有開啟 dev 模式
				key = strings.Join(splited[:2], "-")
			}
		}
		_, out, _ := r.DockerService.GetCreatedTimeOfFirstServiceInStack(line[0])
		out = strings.TrimSuffix(out, "\n")
		t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", out)
		line = append(line, uptime(t))
		stacks[key] = append(stacks[key], line)
	}
	ctx.ViewData("stacks", stacks)
	ctx.View("stack.html")
}

func (r *StackRoutes) DeployStack(ctx iris.Context) {
	d := &models.Deploy{}
	ctx.ReadJSON(d)
	start := time.Now()
	indent, _ := json.MarshalIndent(d, "", " ")

	ctx.StreamWriter(pipe.Printf("Received deploy request: %v", string(indent)))
	err := r.deploy(&ctx, *d)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.StreamWriter(pipe.Printf("Resolving in %v, done.", time.Since(start)))
}

func (r *StackRoutes) RemoveStack(ctx iris.Context) {
	stack := ctx.Params().Get("stack")
	_, _, err := r.DockerStack.RmLike(stack)
	if err != nil {
		ctx.Application().Logger().Warn(err.Error())
		ctx.WriteString(err.Error())
	}
	ctx.Redirect("/")
}

func uptime(t time.Time) string {
	d := time.Since(t)
	return fmt.Sprintf("up %.2f hours", d.Hours())
}

func publishedPort(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func (r *StackRoutes) deploy(ctx *iris.Context, d models.Deploy) error {
	wd := r.Workspace.GetWd(d.CleanUp, d.Project)
	opts := cmd.Options{Ctx: ctx, Pwd: wd.Path}
	d.Dev.PublishPort = d.Dev.Port

	(*ctx).StreamWriter(pipe.Printf("\nDeploying '%v'...\n", d.Yaml))

	gpmDir := "repo"
	group, err := r.Gpm.Install(&opts, gpmDir, &d)
	if err != nil {
		return err
	}
	var yamls []models.Yaml
	repo := path.Join(wd.Path, gpmDir)

	if !group {
		yml := path.Join(repo, "docker-compose.yml")
		dirs, err := collectDirs(repo)
		if err != nil {
			return err
		}
		err = r.GenYaml.Gen(&opts, dirs, yml, &d)
		if err != nil {
			return err
		}
		yamls = append(yamls, models.Yaml{
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
			err := r.GenYaml.Gen(&opts, flat, yml, &d)
			if err != nil {
				return err
			}
			yamls = append(yamls, models.Yaml{
				Group: "",
				Path:  yml,
			})
		} else {
			for group, dirs := range deployGroups {
				yml := path.Join(repo, group, fmt.Sprintf("docker-compose-%v.yml", group))
				err := r.GenYaml.Gen(&opts, dirs, yml, &d)
				if err != nil {
					return err
				}
				yamls = append(yamls, models.Yaml{
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

	err = r.DockerStack.Deploy(&opts, yamls, &d)
	if err != nil {
		return err
	}

	return nil
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
