package hello_world

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestHelloWorldViewIsNotEmpty(t *testing.T) {
	model := NewHelloWorld()

	view := model.View().Content
	if view == "" {
		t.Fatal("hello-world view is empty")
	}

	for _, want := range []string{"Hello, World!", "Press q to quit"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view, got %q", want, view)
		}
	}
}

func TestHelloWorldQuitKeys(t *testing.T) {
	model := NewHelloWorld()

	for _, msg := range []tea.KeyPressMsg{
		tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}),
		tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}),
	} {
		_, cmd := model.Update(msg)
		if cmd == nil {
			t.Fatal("expected quit command")
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("expected quit message, got %T", cmd())
		}
	}
}
