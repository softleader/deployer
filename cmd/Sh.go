package cmd

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	"errors"
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
