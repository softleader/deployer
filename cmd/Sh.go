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

func (sh *Sh) Exec(commands ...string) (string, string, error) {
	arg := strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	cmd.Dir = sh.Ws.Path
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", "", errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}
	return arg, out.String(), nil
}

func (sh *Sh) ExecPipe(ctx *iris.Context, commands ...string) {
	arg := strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	cmd.Dir = sh.Ws.Path

	(*ctx).StreamWriter(pipe.Printf("$ %v\n", arg))

	sw := streamWriter{ctx: ctx}
	cmd.Stdout = sw
	cmd.Stderr = sw
	cmd.Run()
}

type streamWriter struct {
	ctx *iris.Context
}

func (sw streamWriter) Write(bs []byte) (n int, err error) {
	(*sw.ctx).StreamWriter(pipe.Print(string(bs)))
	return len(bs), nil
}
