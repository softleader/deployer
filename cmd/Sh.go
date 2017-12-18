package cmd

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	"errors"
)

type sh struct {
	Wd
}

func Sh() sh {
	return sh{Wd: NewWd()}
}

func (sh sh) Exec(commands ...string) (string, error) {
	arg := strings.Join(commands, " ")
	fmt.Println("$", arg)
	cmd := exec.Command("sh", "-c", arg)
	cmd.Dir = sh.Wd.Path
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}
	return fmt.Sprintf("$ %v\n%v", arg, out.String()), nil
}
