package cmd

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	"errors"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/pipe"
)

type Options struct {
	Ctx   *iris.Context
	Pwd   string
	Debug bool
}

type output struct {
	ctx   *iris.Context
	buf   bytes.Buffer
	Debug bool
}

func Exec(opts *Options, commands ...string) (arg string, out string, err error) {
	arg = strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	if opts.Pwd != "" {
		cmd.Dir = opts.Pwd
	}

	if opts.Ctx != nil {
		(*opts.Ctx).StreamWriter(pipe.Printf("$ %v\n", arg))
	}

	stdout := output{ctx: opts.Ctx, Debug: opts.Debug}
	stderr := output{ctx: opts.Ctx, Debug: opts.Debug}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", "", errors.New(fmt.Sprint(err) + ": " + stderr.buf.String())
	}

	return arg, stdout.buf.String(), nil
}

func (o *output) Write(b []byte) (n int, err error) {
	o.buf.Write(b)
	if o.ctx != nil {
		s := string(b)
		if o.Debug {
			fmt.Print(s)
		}
		(*o.ctx).StreamWriter(pipe.Print(s))
	}
	return len(b), nil
}
