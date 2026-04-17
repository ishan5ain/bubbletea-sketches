package sketches

import (
	"sort"

	tea "charm.land/bubbletea/v2"

	directory_selector "github.com/ishansain/bubbletea-sketches/internal/sketches/directory_selector"
	flexible_key_value_pair_list "github.com/ishansain/bubbletea-sketches/internal/sketches/flexible_key_value_pair_list"
	hello_world "github.com/ishansain/bubbletea-sketches/internal/sketches/hello_world"
	styled_hello_world "github.com/ishansain/bubbletea-sketches/internal/sketches/styled_hello_world"
)

type Factory func() tea.Model

const defaultSketchName = "hello-world"

var registry = map[string]Factory{
	defaultSketchName:              hello_world.NewHelloWorld,
	"directory-selector":           directory_selector.NewDirectorySelector,
	"flexible-key-value-pair-list": flexible_key_value_pair_list.NewFlexibleKeyValuePairList,
	"styled-hello-world":           styled_hello_world.NewStyledHelloWorld,
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
