package command

import (
	"bytes"
	"log"
	"os/exec"

	"golang.org/x/term"
)

func ExecCmd(cmd ...string) (string, error) {
	res := exec.Command(cmd[0], cmd[1:]...)
	var buf bytes.Buffer
	res.Stdout = &buf
	err := res.Run()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetGitStatus() string {
	res, err := ExecCmd("git", "status", "-s")
	if err != nil {
		return ""
	}

	return res
}

func GitAddAll() (string, error) {
	return ExecCmd("git", "add", ".")
}

func GetStagedFiles() string {
	res, err := ExecCmd("git", "diff", "--staged")
	if err != nil {
		log.Panic(err.Error())
	}
	return res
}

func GetCmdWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		return 20
	}
	return width
}
