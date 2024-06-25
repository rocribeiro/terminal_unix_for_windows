package console

import (
	"fmt"
	"os"
	"path/filepath"
	"terminal_go/internal/autocomplete"
	"terminal_go/internal/comandos"

	"github.com/fatih/color"
	"github.com/peterh/liner"
)

func CreateConsole() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	// Funções de cor
	termColor := color.New(color.FgMagenta).SprintFunc() // Cor exclusiva para "GoTerm"
	pathColor := color.New(color.FgCyan).SprintFunc()    // Cor para o caminho

	// Obter o diretório de trabalho atual
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		dir = "unknown"
	}

	// Configurar o prompt para mostrar o caminho atual com cores
	coloredPrompt := fmt.Sprintf("%s %s> ", termColor("GoTerm-For-Windows"), pathColor(filepath.Clean(dir)))

	// Defina o prompt fixo uma vez
	// line.SetPrompt(coloredPrompt)

	// Configura a função de autocompletação
	line.SetCompleter(func(input string) (c []string) {
		return autocomplete.Complete(input)
	})

	for {
		// Display the colored prompt to the user
		fmt.Print(coloredPrompt)

		// Read the input using an empty prompt string
		input, err := line.Prompt("")
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
