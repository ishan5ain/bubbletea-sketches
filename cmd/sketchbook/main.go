package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/ishansain/bubbletea-sketches/internal/sketches"
)

func main() {
	model, err := selectSketch(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	program := tea.NewProgram(model, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run sketch: %v\n", err)
		os.Exit(1)
	}
}

func selectSketch(args []string) (tea.Model, error) {
	name := sketches.DefaultName()
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		name = args[0]
	}

	factory, ok := sketches.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown sketch %q; available sketches: %s", name, strings.Join(sketches.Names(), ", "))
	}

	return factory(), nil
}
