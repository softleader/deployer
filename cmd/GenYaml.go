package cmd

import (
	"github.com/softleader/deployer/models"
	"strings"
	"github.com/softleader/deployer/pipe"
	"regexp"
	"strconv"
	"os"
	"fmt"
	"log"
)

type GenYaml struct {
	cmd string
}

func NewGenYaml(cmd string) *GenYaml {
	if cmd == "" {
		cmd = "gen-yaml"
	}
	genYaml := GenYaml{cmd: cmd}
	cmd, out, err := genYaml.Version()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("  $ %v: %v", cmd, out)
	return &genYaml
}

func (gy *GenYaml) Version() (arg string, out string, err error) {
	return Exec(&Options{}, gy.cmd, "--version")
}

func (gy *GenYaml) Gen(opts *Options, dirs []string, output string, d *models.Deploy) error {
	_, out, err := gen(gy.cmd, opts, output, d, strings.Join(dirs, " "))
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

func gen(cmd string, opts *Options, output string, d *models.Deploy, dirs ...string) (arg string, out string, err error) {
	commands := []string{cmd, "-s", d.Style, "-o", output}
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

	return Exec(opts, commands...)
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
