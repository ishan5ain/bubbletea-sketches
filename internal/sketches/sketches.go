package sketches

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type Factory func() tea.Model

const defaultSketchName = "hello-world"

var registry = map[string]Factory{
	defaultSketchName: NewHelloWorld,
}

func DefaultName() string {
	return defaultSketchName
}

func Get(name string) (Factory, bool) {
	factory, ok := registry[name]
	return factory, ok
}

func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}
