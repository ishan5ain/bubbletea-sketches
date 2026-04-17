package main

import (
	"fmt"
	"os"

	"github.com/ishansain/bubbletea-sketches/internal/sketchbookapp"
)

func main() {
	if err := sketchbookapp.RunBrowser(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
