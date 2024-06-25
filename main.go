package main

import (
	"fmt"
	"os"

	"terminal_go/internal/autocomplete"
	"terminal_go/internal/comandos"
	"terminal_go/internal/console"

	"github.com/peterh/liner"
)

func main() {
	console.CreateConsole()
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		return autocomplete.Complete(line)
	})

	for {
		input, err := line.Prompt("> ")
		if err == liner.ErrPromptAborted {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		line.AppendHistory(input)

		if err = comandos.ExecInput(input); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
