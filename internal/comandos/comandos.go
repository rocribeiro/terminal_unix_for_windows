package comandos

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var errNoPath = errors.New("path required")

func ExecInput(input string) error {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return errNoPath
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	cmdName, cmdArgs := convertCommand(args[0], args[1:]...)
	cmd := exec.Command(cmdName, cmdArgs...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func convertCommand(cmd string, args ...string) (string, []string) {
	if runtime.GOOS == "windows" {
		switch cmd {
		case "ls":
			return "cmd", append([]string{"/C", "dir"}, args...)
		case "cat":
			return "cmd", append([]string{"/C", "type"}, args...)
		case "rm":
			return "cmd", append([]string{"/C", "del"}, args...)
		case "cp":
			return "cmd", append([]string{"/C", "copy"}, args...)
		case "mv":
			return "cmd", append([]string{"/C", "move"}, args...)
		case "clear":
			return "cmd", []string{"/C", "cls"}
		case "pwd":
			return "cmd", []string{"/C", "cd"}
		}
	}
	return cmd, args
}
