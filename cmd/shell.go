package cmd

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	"github.com/kataras/iris"
	"github.com/softleader/deployer/pipe"
	"gopkg.in/yaml.v2"
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

type ExecError struct {
	Cmd    string
	Err    string // underlying error
	Stdout string
	Stderr string
}

func (e *ExecError) Error() string {
	b, _ := yaml.Marshal(e)
	return string(b)
}

func NewExecError(cmd string, err error, stdout string, stderr string) (e *ExecError) {
	e = &ExecError{
		Cmd:    cmd,
		Stdout: stdout,
		Stderr: stderr,
	}
	if err != nil {
		e.Err = err.Error()
	}
	return
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
	if err != nil || stderr.buf.Len() > 0 {
		return "", "", NewExecError(arg, err, stdout.buf.String(), stderr.buf.String())
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
