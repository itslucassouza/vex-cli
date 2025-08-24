package main

import (
	"fmt"
	"os"

	tui "github.com/itslucassouza/vex-cli/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
