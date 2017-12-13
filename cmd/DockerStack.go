package cmd

import (
	"os/exec"
	"os"
)

type DockerStack struct{}

func (DockerStack) Ls() string {
	out, err := exec.Command("sh", "-c", "docker stack ls").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return string(out)
}

func (DockerStack) Rm(stack string) string {
	out, err := exec.Command("sh", "-c", "docker stack rm "+stack).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return string(out)
}

func (DockerStack) Deploy(stack string, file string) string {
	out, err := exec.Command("sh", "-c", "docker stack deploy -c "+file+" "+stack).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return string(out)
}
