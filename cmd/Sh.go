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

type Sh struct {
	Ws
}

func NewSh(wd Ws) *Sh {
	return &Sh{Ws: wd}
}

type out struct {
	ctx *iris.Context
	buf bytes.Buffer
}

func (sh *Sh) Exec(ctx *iris.Context, commands ...string) (string, string, error) {
	arg := strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	cmd.Dir = sh.Ws.Path

	if ctx != nil {
		(*ctx).StreamWriter(pipe.Printf("$ %v\n", arg))
	}

	stdout := out{ctx: ctx}
	stderr := out{ctx: ctx}
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
