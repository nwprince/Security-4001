package node

import (
	"fmt"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

// ExecuteCommand will run an arbitrary command
func ExecuteCommand(name string, arg string) (stdout string, stderr string) {
	var cmd *exec.Cmd

	argS, err := shellwords.Parse(arg)
	if err != nil {
		return "", fmt.Sprintf("there was an err with args: %2\r\n%s", arg, err.Error())
	}

	cmd = exec.Command(name, argS...)

	out, err := cmd.CombinedOutput()
	stdout = string(out)
	stderr = ""

	if err != nil {
		stderr = err.Error()
	}
	return stdout, stderr
}
