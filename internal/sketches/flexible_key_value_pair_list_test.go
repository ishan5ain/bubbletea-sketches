package sketches

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func TestAutocompleteOptionalKeySingleMatch(t *testing.T) {
	got := autocompleteOptionalKey("tim", kvDefaultSchema().optional, nil)
	if got.input != "timeout" {
		t.Fatalf("expected %q, got %q", "timeout", got.input)
	}
	if !got.advanced {
		t.Fatal("expected autocomplete to advance for single match")
	}
}

func TestAutocompleteOptionalKeyExcludesUsedKeys(t *testing.T) {
	got := autocompleteOptionalKey("re", kvDefaultSchema().optional, []string{"region"})
	if len(got.candidates) != 1 || got.candidates[0] != "retries" {
		t.Fatalf("unexpected candidates: %#v", got.candidates)
	}
}

func TestAutocompleteOptionalKeyUsesLongestCommonPrefix(t *testing.T) {
	got := autocompleteOptionalKey("r", []string{"region", "retries"}, nil)
	if got.input != "re" {
		t.Fatalf("expected %q, got %q", "re", got.input)
	}
	if !got.advanced {
		t.Fatal("expected autocomplete to advance to common prefix")
	}
}

func TestAutocompleteOptionalKeyNoMatch(t *testing.T) {
	got := autocompleteOptionalKey("zzz", kvDefaultSchema().optional, nil)
	if got.input != "zzz" {
		t.Fatalf("expected input unchanged, got %q", got.input)
	}
	if len(got.candidates) != 0 {
		t.Fatalf("expected no candidates, got %#v", got.candidates)
	}
}

func TestFlexibleKeyValuePairListTypingRequiredValue(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Text: "/", Code: '/'}))
	got := next.(flexibleKeyValuePairListModel)
	if got.required[0].value != "/" {
		t.Fatalf("expected required value %q, got %q", "/", got.required[0].value)
	}
}

func TestFlexibleKeyValuePairListArrowNavigation(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	got := next.(flexibleKeyValuePairListModel)
	if got.focusRow != 1 || got.focusField != kvValueField {
		t.Fatalf("expected focus on second required value field, got row=%d field=%d", got.focusRow, got.focusField)
	}

	next, _ = got.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	got = next.(flexibleKeyValuePairListModel)
	if got.focusRow != 2 || got.focusField != kvKeyField {
		t.Fatalf("expected focus on optional key field, got row=%d field=%d", got.focusRow, got.focusField)
	}
}

func TestFlexibleKeyValuePairListTabAutocompletesOptionalKey(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.focusRow = len(model.required)
	model.focusField = kvKeyField
	model.optional[0].key = "tim"
	model.recalculate()

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	got := next.(flexibleKeyValuePairListModel)
	if got.optional[0].key != "timeout" {
		t.Fatalf("expected autocomplete to fill timeout, got %q", got.optional[0].key)
	}
}

func TestFlexibleKeyValuePairListTabAdvancesFocus(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.focusRow = len(model.required)
	model.focusField = kvKeyField
	model.optional[0].key = "timeout"
	model.recalculate()

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	got := next.(flexibleKeyValuePairListModel)
	if got.focusField != kvValueField {
		t.Fatalf("expected focus to advance to value field, got %d", got.focusField)
	}
}

func TestFlexibleKeyValuePairListKeepsTrailingBlankOptionalRow(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.focusRow = len(model.required)
	model.focusField = kvKeyField
	model.optional[0].key = "timeout"
	model.optional[0].value = "30s"
	model.recalculate()

	if len(model.optional) != 2 {
		t.Fatalf("expected a new trailing blank row, got %d optional rows", len(model.optional))
	}
	if model.optional[1].key != "" || model.optional[1].value != "" {
		t.Fatalf("expected trailing blank row, got %#v", model.optional[1])
	}
}

func TestFlexibleKeyValuePairListDuplicateKeysValidate(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.optional[0] = kvRow{key: "timeout", value: "30s"}
	model.optional = append(model.optional, kvRow{key: "timeout", value: "60s"})
	model.recalculate()

	if model.optional[0].err != "duplicate key" || model.optional[1].err != "duplicate key" {
		t.Fatalf("expected duplicate key errors, got %#v %#v", model.optional[0].err, model.optional[1].err)
	}
}

func TestFlexibleKeyValuePairListQuitKeys(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()

	for _, msg := range []tea.KeyPressMsg{
		tea.KeyPressMsg(tea.Key{Code: 'q', Mod: tea.ModCtrl}),
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

func TestFlexibleKeyValuePairListPlainQTypesIntoFocusedField(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}))
	if cmd != nil {
		t.Fatal("expected plain q to type into the field, not quit")
	}

	got := next.(flexibleKeyValuePairListModel)
	if got.required[0].value != "q" {
		t.Fatalf("expected plain q to append to the focused field, got %q", got.required[0].value)
	}
}

func TestFlexibleKeyValuePairListViewContainsPreviewAndSuggestions(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.focusRow = len(model.required)
	model.focusField = kvKeyField
	model.optional[0].key = "re"
	model.recalculate()

	view := model.View().Content
	for _, want := range []string{"Preview", "suggestions:", "key", "value"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q in view, got %q", want, view)
		}
	}
}

func TestFlexibleKeyValuePairListRequiredValidationIsInline(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	view := model.View().Content

	if !strings.Contains(view, "required value") {
		t.Fatalf("expected inline required validation in view, got %q", view)
	}
	if strings.Contains(view, "! required value") {
		t.Fatalf("expected no legacy second-line error marker, got %q", view)
	}
}

func TestFlexibleKeyValuePairListUnknownKeyUsesBadgeStyle(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.optional[0] = kvRow{key: "mystery", value: "123"}
	model.recalculate()

	view := model.View().Content
	want := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B91C1C")).
		Background(lipgloss.Color("#FEE2E2")).
		Padding(0, 1).
		Render("unknown key")

	if !strings.Contains(view, want) {
		t.Fatalf("expected styled unknown-key badge in view, got %q", view)
	}
}

func TestFlexibleKeyValuePairListEmptyFieldsShowGhostLabels(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	view := model.View().Content

	for _, want := range []string{"key", "value"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected ghost label %q in view, got %q", want, view)
		}
	}
}

func TestFlexibleKeyValuePairListFocusedCursorFollowsTypedContent(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.focusRow = len(model.required)
	model.focusField = kvKeyField
	model.optional[0].key = "tim"
	model.suggestion = "eout"

	view := model.View().Content
	typed := lipgloss.NewStyle().Foreground(lipgloss.Color("#DBEAFE")).Bold(true).Render("tim")
	cursor := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#F8FAFC")).Render(" ")
	suggestion := lipgloss.NewStyle().Foreground(lipgloss.Color("#00eeff")).Background(lipgloss.Color("#b0b0b0")).Render("eout")

	typedIdx := strings.Index(view, typed)
	cursorIdx := strings.Index(view, cursor)
	suggestionIdx := strings.Index(view, suggestion)
	if typedIdx == -1 || cursorIdx == -1 || suggestionIdx == -1 {
		t.Fatalf("expected typed text, cursor, and suggestion in view, got %q", view)
	}
	if !(typedIdx < cursorIdx && cursorIdx < suggestionIdx) {
		t.Fatalf("expected cursor between typed text and suggestion, got %q", view)
	}
}

func TestFlexibleKeyValuePairListCompletedRequiredBadgeUsesBlue(t *testing.T) {
	model := newFlexibleKeyValuePairListModelForTest()
	model.required[0].value = "/users"
	model.recalculate()

	view := model.View().Content
	want := lipgloss.NewStyle().Foreground(lipgloss.Color("#92b0d2")).Render("required")
	if !strings.Contains(view, want) {
		t.Fatalf("expected completed required badge style in view, got %q", view)
	}
}

func newFlexibleKeyValuePairListModelForTest() flexibleKeyValuePairListModel {
	model := NewFlexibleKeyValuePairList().(flexibleKeyValuePairListModel)
	return model
}
