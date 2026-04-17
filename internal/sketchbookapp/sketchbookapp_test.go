package sketchbookapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestPageItemsPrefixesItemsPerPage(t *testing.T) {
	names := []string{
		"alpha", "beta", "gamma", "delta", "epsilon",
		"zeta", "eta", "theta", "iota", "kappa", "lambda",
	}

	items := pageItems(names, 0)
	if len(items) != 10 {
		t.Fatalf("expected 10 items on first page, got %d", len(items))
	}

	first, ok := items[0].(sketchItem)
	if !ok {
		t.Fatalf("expected sketchItem, got %T", items[0])
	}
	if first.Title() != "0 alpha" {
		t.Fatalf("expected first prefix to be 0, got %q", first.Title())
	}

	last, ok := items[9].(sketchItem)
	if !ok {
		t.Fatalf("expected sketchItem, got %T", items[9])
	}
	if last.Title() != "9 kappa" {
		t.Fatalf("expected last prefix to be 9, got %q", last.Title())
	}

	tail := pageItems(names, 1)
	if len(tail) != 1 {
		t.Fatalf("expected 1 item on second page, got %d", len(tail))
	}

	tailItem := tail[0].(sketchItem)
	if tailItem.Title() != "0 lambda" {
		t.Fatalf("expected second page prefix reset, got %q", tailItem.Title())
	}
}

func TestBrowserPagesClampAndSelectDigits(t *testing.T) {
	model := newBrowserModel([]string{
		"s0", "s1", "s2", "s3", "s4",
		"s5", "s6", "s7", "s8", "s9",
		"s10", "s11",
	})

	if model.page != 0 {
		t.Fatalf("expected initial page 0, got %d", model.page)
	}

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyLeft}))
	model = next.(browserModel)
	if model.page != 0 {
		t.Fatalf("expected left at first page to clamp, got %d", model.page)
	}

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyRight}))
	model = next.(browserModel)
	if model.page != 1 {
		t.Fatalf("expected right to advance page, got %d", model.page)
	}

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyRight}))
	model = next.(browserModel)
	if model.page != 1 {
		t.Fatalf("expected right at last page to clamp, got %d", model.page)
	}

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Text: "0", Code: '0'}))
	model = next.(browserModel)
	if cmd == nil {
		t.Fatal("expected digit selection to quit the browser")
	}
	if model.selected != "s10" {
		t.Fatalf("expected page-local digit 0 to select s10, got %q", model.selected)
	}
}

func TestBrowserViewShowsControlsAndPageTitle(t *testing.T) {
	model := newBrowserModel([]string{"hello-world", "directory-selector"})
	view := model.View().Content

	for _, want := range []string{"Sketchbook [1/1]", "0-9 run", "left/right page"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view, got %q", want, view)
		}
	}
}

func TestBrowserViewShowsItemsAfterWindowSizing(t *testing.T) {
	model := newBrowserModel([]string{
		"hello-world",
		"directory-selector",
		"flexible-key-value-pair-list",
		"styled-hello-world",
	})

	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	model = next.(browserModel)

	view := model.View().Content
	for _, want := range []string{"0 hello-world", "1 directory-selector"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view after sizing, got %q", want, view)
		}
	}
}

func TestBrowserEnterRunsHighlightedSketch(t *testing.T) {
	model := newBrowserModel([]string{
		"hello-world",
		"directory-selector",
		"flexible-key-value-pair-list",
		"styled-hello-world",
	})

	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	model = next.(browserModel)

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(browserModel)

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(browserModel)

	if cmd == nil {
		t.Fatal("expected enter to quit the browser")
	}
	if model.selected != "directory-selector" {
		t.Fatalf("expected enter to select the highlighted sketch, got %q", model.selected)
	}
}
