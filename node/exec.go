package node

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

// ExecuteCommand will run an arbitrary command
func ExecuteCommand(name string, arg string) (stdout string, stderr string) {
	var cmd *exec.Cmd
	var argS []string
	var err error
	isScript := false

	if strings.Contains(arg, ">>") {
		argS = strings.Fields(arg)
		isScript = true
	} else {
		argS, err = shellwords.Parse(arg)
	}
	if err != nil {
		return "", fmt.Sprintf("there was an err with args: %s\r\n%s", arg, err.Error())
	}

	if isScript {
		scriptCmd := argS[0 : len(argS)-2]
		strCmd := strings.Join(scriptCmd, " ")
		strCmd = strCmd + "\n"
		f, err := os.OpenFile("exec.sh", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return "", err.Error()
		}
		_, err = f.WriteString(strCmd)
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
		return "", ""
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
