package comandos

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var errNoPath = errors.New("path required")

func ExecInput(input string) error {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			// Se o caminho não for especificado, mostrar o diretório atual
			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Error: %v", err)
			}
			fmt.Println(dir)
			return nil
		}
		// Se o caminho for especificado, mudar para o diretório especificado
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "help":
		printHelp()
		return nil
	case "ls":
		return listFiles()
	default:
		cmdName, cmdArgs := convertCommand(args[0], args[1:]...)
		// Verifique se o comando está disponível no PATH
		if _, err := exec.LookPath(cmdName); err != nil {
			return fmt.Errorf("\"%s\": comando não reconhecido", cmdName)
		}
		cmd := exec.Command(cmdName, cmdArgs...)

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	}

	return nil
}

func convertCommand(cmd string, args ...string) (string, []string) {
	if runtime.GOOS == "windows" {
		switch cmd {
		case "ls":
			return "powershell", append([]string{"-Command", "Get-ChildItem"}, args...)
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

func printHelp() {
	fmt.Println("Comandos suportados:")
	fmt.Println("  cd [path]   - Muda o diretório atual para [path] ou mostra o diretório atual se nenhum path for especificado")
	fmt.Println("  exit        - Sai do programa")
	fmt.Println("  help        - Mostra esta mensagem de ajuda")
	fmt.Println("  ls          - Lista arquivos e diretórios")
	fmt.Println("  cat         - Mostra o conteúdo de um arquivo (adaptado para Windows)")
	fmt.Println("  rm          - Remove arquivos (adaptado para Windows)")
	fmt.Println("  cp          - Copia arquivos (adaptado para Windows)")
	fmt.Println("  mv          - Move arquivos (adaptado para Windows)")
	fmt.Println("  clear       - Limpa a tela (adaptado para Windows)")
	fmt.Println("  pwd         - Mostra o diretório atual (adaptado para Windows)")
	fmt.Println("Outros comandos serão executados como estão, se disponíveis no PATH.")
}

func listFiles() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	// Exibir o cabeçalho
	fmt.Println("Mode                 LastWriteTime         Length Name")
	fmt.Println("----                 -------------         ------ ----")

	// Exibir arquivos no diretório com cores diferentes
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return nil
		}

		var colorFunc func(format string, a ...interface{}) string
		switch {
		case info.IsDir():
			colorFunc = color.New(color.FgBlue).SprintfFunc()
		case strings.HasSuffix(info.Name(), ".go"):
			colorFunc = color.New(color.FgGreen).SprintfFunc()
		case strings.HasSuffix(info.Name(), ".txt"):
			colorFunc = color.New(color.FgYellow).SprintfFunc()
		default:
			colorFunc = color.New(color.FgWhite).SprintfFunc()
		}

		mode := "d-----"
		if !info.IsDir() {
			mode = "-a----"
		}
		lastWriteTime := info.ModTime().Format("02/01/2006 15:04")
		length := ""
		if !info.IsDir() {
			length = fmt.Sprintf("%d", info.Size())
		}
		name := filepath.Base(path)

		fmt.Println(colorFunc("%-20s %-20s %-6s %s", mode, lastWriteTime, length, name))
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}
