package directory_selector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestAnalyzeDirectoryInputRelativeSingleMatch(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "docs")
	mustMkdir(t, root, "downloads")

	got := analyzeDirectoryInput("doc", root, root)
	if got.Suggestion != "s/" {
		t.Fatalf("expected single-match suggestion %q, got %q", "s/", got.Suggestion)
	}
	if len(got.Matches) != 1 || got.Matches[0] != "docs/" {
		t.Fatalf("unexpected matches: %#v", got.Matches)
	}
}

func TestCompleteDirectoryInputAbsoluteSingleMatch(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "projects")

	input := filepath.Join(root, "pro")
	got := completeDirectoryInput(input, root, root)
	want := filepath.Join(root, "projects") + string(filepath.Separator)
	if got.Input != want {
		t.Fatalf("expected %q, got %q", want, got.Input)
	}
}

func TestAnalyzeDirectoryInputExpandsTilde(t *testing.T) {
	home := t.TempDir()
	mustMkdir(t, home, "workspace")

	got := analyzeDirectoryInput("~/wor", home, home)
	if len(got.Matches) != 1 || got.Matches[0] != "~/workspace/" {
		t.Fatalf("unexpected matches: %#v", got.Matches)
	}
}

func TestCompleteDirectoryInputUsesLongestCommonPrefix(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "projects-api")
	mustMkdir(t, root, "projects-web")

	got := completeDirectoryInput("pro", root, root)
	if got.Input != "projects-" {
		t.Fatalf("expected longest common prefix completion, got %q", got.Input)
	}
	if len(got.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got.Matches))
	}
}

func TestAnalyzeDirectoryInputNoMatches(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "docs")

	got := analyzeDirectoryInput("zzz", root, root)
	if got.Suggestion != "" {
		t.Fatalf("expected no suggestion, got %q", got.Suggestion)
	}
	if len(got.Matches) != 0 {
		t.Fatalf("expected no matches, got %#v", got.Matches)
	}
}

func TestAnalyzeDirectoryInputExcludesFiles(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "src")
	filePath := filepath.Join(root, "script.sh")
	if err := os.WriteFile(filePath, []byte("echo hi"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := analyzeDirectoryInput("s", root, root)
	if len(got.Matches) != 1 || got.Matches[0] != "src/" {
		t.Fatalf("unexpected matches: %#v", got.Matches)
	}
}

func TestDirectorySelectorTypingAndBackspace(t *testing.T) {
	root := t.TempDir()
	model := newDirectorySelectorModelForTest(root, root)

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Text: "d", Code: 'd'}))
	typed := next.(directorySelectorModel)
	if typed.input != "d" {
		t.Fatalf("expected input %q, got %q", "d", typed.input)
	}

	backspaced, _ := typed.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyBackspace}))
	if got := backspaced.(directorySelectorModel).input; got != "" {
		t.Fatalf("expected empty input after backspace, got %q", got)
	}
}

func TestDirectorySelectorTabUpdatesInput(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "documents")
	model := newDirectorySelectorModelForTest(root, root)
	model.input = "doc"
	model.refreshCompletion()

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	if got := next.(directorySelectorModel).input; got != "documents/" {
		t.Fatalf("expected completed input, got %q", got)
	}
}

func TestDirectorySelectorEnterSelectsValidDirectory(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, root, "documents")
	model := newDirectorySelectorModelForTest(root, root)
	model.input = "documents"

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	got := next.(directorySelectorModel)
	want := filepath.Join(root, "documents")
	if got.selectedPath != want {
		t.Fatalf("expected selected path %q, got %q", want, got.selectedPath)
	}
}

func TestDirectorySelectorEnterShowsErrorForInvalidPath(t *testing.T) {
	root := t.TempDir()
	model := newDirectorySelectorModelForTest(root, root)
	model.input = "missing"

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	got := next.(directorySelectorModel)
	if got.errorMsg == "" {
		t.Fatal("expected an error message for an invalid path")
	}
}

func TestDirectorySelectorQuitKeys(t *testing.T) {
	root := t.TempDir()
	model := newDirectorySelectorModelForTest(root, root)

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

func TestDirectorySelectorViewContainsPromptAndSelectedState(t *testing.T) {
	root := t.TempDir()
	model := newDirectorySelectorModelForTest(root, root)
	model.input = "doc"
	model.suggestion = "uments/"
	model.selectedPath = filepath.Join(root, "documents")

	view := model.View().Content
	for _, want := range []string{"directory-selector", "cd ", "Selected:", "tab autocomplete"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view, got %q", want, view)
		}
	}
}

func newDirectorySelectorModelForTest(cwd string, home string) directorySelectorModel {
	model := directorySelectorModel{
		cwd:  cwd,
		home: home,
	}
	model.refreshCompletion()
	return model
}

func mustMkdir(t *testing.T, root string, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, name), 0o755); err != nil {
		t.Fatalf("mkdir %q: %v", name, err)
	}
}
