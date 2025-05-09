package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func PromptInput(scanner *bufio.Scanner, label string) string {
	fmt.Print(label + ": ")
	// wait for user input till the key enter is pressed
	// then the input will stored internaly by scanner
	scanner.Scan()
	// get input from scanner
	return strings.TrimSpace(scanner.Text())
}

func PromptRequiredInput(scanner *bufio.Scanner, label string) string {
	for {
		input := PromptInput(scanner, label)
		if input != "" {
			return input
		}
		SetRespond(label + " is required", "error")
	}
}

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default: // linux, macOS, etc.
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}