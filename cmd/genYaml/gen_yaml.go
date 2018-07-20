package genYaml

import (
	"github.com/softleader/deployer/models"
	"strings"
	"regexp"
	"strconv"
	"github.com/kataras/iris/core/errors"
	"github.com/softleader/deployer/cmd"
)

var Cmd string

func Version() (arg string, out string, err error) {
	return cmd.Exec(&cmd.Options{}, Cmd, "--version")
}

func Gen(opts *cmd.Options, dirs []string, output string, d *models.Deploy) error {
	if len(dirs) <= 0 {
		return errors.New("dirs is required")
	}
	_, out, err := gen(opts, output, d, strings.Join(dirs, " "))
	if err != nil {
		return err
	}
	err = updateDevPort(out, d)
	if err != nil {
		return err
	}
	return nil
}

func gen(opts *cmd.Options, output string, d *models.Deploy, dirs ...string) (arg string, out string, err error) {
	commands := []string{Cmd, "-s", d.Style, "-o", output}
	if d.Silently {
		commands = append(commands, "-S")
	}
	if d.Net0 != "" {
		commands = append(commands, "--net0", d.Net0)
	}
	if d.Volume0 != "" {
		commands = append(commands, "--volume0", d.Volume0)
	}
	if d.Dev.IpAddress != "" {
		commands = append(commands, d.Dev.String())
	}
	commands = append(commands, strings.Join(dirs, " "))

	return cmd.Exec(opts, commands...)
}

func updateDevPort(out string, d *models.Deploy) error {
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
