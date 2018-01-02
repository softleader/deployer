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

type Shell struct {
}

type Options struct {
	Ctx *iris.Context
	Pwd string
}

func NewShell() *Shell {
	return &Shell{}
}

type out struct {
	ctx *iris.Context
	buf bytes.Buffer
}

func (sh *Shell) Exec(opts *Options, commands ...string) (string, string, error) {
	arg := strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	if opts.Pwd != "" {
		cmd.Dir = opts.Pwd
	}

	if opts.Ctx != nil {
		(*opts.Ctx).StreamWriter(pipe.Printf("$ %v\n", arg))
	}

	stdout := out{ctx: opts.Ctx}
	stderr := out{ctx: opts.Ctx}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", "", errors.New(fmt.Sprint(err) + ": " + stderr.buf.String())
	}

	return arg, stdout.buf.String(), nil
}

func (o *out) Write(b []byte) (n int, err error) {
	o.buf.Write(b)
	if o.ctx != nil {
		(*o.ctx).StreamWriter(pipe.Print(string(b)))
	}
	return len(b), nil
}
