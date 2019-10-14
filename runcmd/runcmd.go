package runcmd

import (
	"bytes"
	"os/exec"
	"strings"
)

func RunCmd(cmd string) (error, string, string) {
	shebang := "bash"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	run := exec.Command(shebang, "-c", cmd)
	run.Stdout = &stdout
	run.Stderr = &stderr
	err := run.Run()
	return err, strings.Trim(stdout.String(), "\n"), strings.Trim(stderr.String(), "\n")
}
