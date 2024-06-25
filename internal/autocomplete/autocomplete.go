package autocomplete

import (
	"os"
	"path/filepath"
	"strings"
)

var commands = []string{"cd", "exit", "ls", "cat", "rm", "cp", "mv", "clear", "pwd", "help"}

func Complete(line string) (c []string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		for _, cmd := range commands {
			if strings.HasPrefix(cmd, parts[0]) {
				c = append(c, cmd)
			}
		}
	} else {
		dir := "."
		prefix := parts[len(parts)-1]
		if strings.Contains(prefix, string(os.PathSeparator)) {
			dir = filepath.Dir(prefix)
			prefix = filepath.Base(prefix)
		}

		matches, err := filepath.Glob(filepath.Join(dir, prefix+"*"))
		if err == nil {
			for _, match := range matches {
				if info, err := os.Stat(match); err == nil && info.IsDir() {
					match += string(os.PathSeparator)
				}
				// Preserve the command part and append the match
				c = append(c, strings.Join(parts[:len(parts)-1], " ")+" "+match)
			}
		}
	}
	return
}
