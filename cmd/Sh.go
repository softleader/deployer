package cmd

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	"errors"
	"github.com/kataras/iris"
	"io"
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

	(*ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprintf(w, "$ %v\n", arg)
		return false
	})
	cmd.Stdout = pipe{ctx: ctx}
	cmd.Stderr = pipe{ctx: ctx}
	cmd.Run()
}

type pipe struct {
	ctx *iris.Context
}

func (p pipe) Write(bs []byte) (n int, err error) {
	(*p.ctx).StreamWriter(func(w io.Writer) bool {
		fmt.Fprint(w, string(bs))
		return false
	})
	return len(bs), nil
}
