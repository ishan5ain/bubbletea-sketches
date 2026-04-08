package sketches

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestStyledHelloWorldViewIsNotEmpty(t *testing.T) {
	model := NewStyledHelloWorld()
	if got := model.View().Content; got == "" {
		t.Fatal("styled-hello-world view is empty")
	}
}

func TestStyledHelloWorldNavigationChangesSections(t *testing.T) {
	model := NewStyledHelloWorld()
	initial := model.View().Content

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyRight}))
	if cmd != nil {
		t.Fatal("expected no command when moving to the next section")
	}

	if got := next.View().Content; got == initial {
		t.Fatal("expected section content to change after right navigation")
	}
}

func TestStyledHelloWorldNavigationWraps(t *testing.T) {
	model := NewStyledHelloWorld()
	wrapped, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyLeft}))

	if got := wrapped.View().Content; !strings.Contains(got, "Composition") {
		t.Fatalf("expected left navigation from first section to wrap to the last section, got %q", got)
	}
}

func TestStyledHelloWorldQuitKeys(t *testing.T) {
	model := NewStyledHelloWorld()

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

func TestStyledHelloWorldViewContainsRepresentativeSections(t *testing.T) {
	model := NewStyledHelloWorld()
	view := model.View().Content

	for _, want := range []string{"Text Emphasis", "Hello, World!", "left/h previous"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view, got %q", want, view)
		}
	}
}
