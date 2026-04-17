package flexible_key_value_pair_list

import (
	"encoding/json"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type kvField int

const (
	kvKeyField kvField = iota
	kvValueField
)

const (
	kvKeyFieldWidth   = 22
	kvValueFieldWidth = 34
)

type kvSchema struct {
	required []string
	optional []string
}

type kvRow struct {
	key      string
	value    string
	required bool
	err      string
}

type kvAutocomplete struct {
	input      string
	suggestion string
	candidates []string
	advanced   bool
}

type flexibleKeyValuePairListModel struct {
	schema      kvSchema
	required    []kvRow
	optional    []kvRow
	focusRow    int
	focusField  kvField
	suggestion  string
	candidates  []string
	preview     string
	previewNote string
}

func NewFlexibleKeyValuePairList() tea.Model {
	schema := kvDefaultSchema()
	required := make([]kvRow, 0, len(schema.required))
	for _, key := range schema.required {
		required = append(required, kvRow{key: key, required: true})
	}

	model := flexibleKeyValuePairListModel{
		schema:     schema,
		required:   required,
		optional:   []kvRow{{}},
		focusRow:   0,
		focusField: kvValueField,
	}
	model.recalculate()
	return model
}

func (m flexibleKeyValuePairListModel) Init() tea.Cmd {
	return nil
}

func (m flexibleKeyValuePairListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		case "up":
			m.moveRow(-1)
		case "down":
			m.moveRow(1)
		case "left":
			m.moveField(-1)
		case "right":
			m.moveField(1)
		case "backspace":
			m.backspace()
		case "tab":
			m.handleTab()
		case "enter":
			m.advanceFocus()
		default:
			if text := msg.Key().Text; text != "" {
				m.appendText(text)
			}
		}
	}

	return m, nil
}

func (m flexibleKeyValuePairListModel) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Background(lipgloss.Color("#111827")).
		Padding(0, 0).
		Render("flexible-key-value-pair-list")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")).
		Render("build REST params or CLI flags with required fields, optional keys, and schema-backed autocomplete")

	lines := []string{title, subtitle, ""}
	for i := 0; i < m.totalRows(); i++ {
		lines = append(lines, m.renderRow(i))
	}

	lines = append(lines, "")
	lines = append(lines, m.renderPreview())
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Render("[arrows] move focus  [tab] autocomplete/advance  [enter] next field  [^q] quit"))

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(1, 2).
		Render(strings.Join(lines, "\n"))

	return tea.NewView(panel)
}

func (m *flexibleKeyValuePairListModel) moveRow(delta int) {
	wasRequired := m.isRequiredRow(m.focusRow)
	next := m.focusRow + delta
	if next < 0 {
		next = 0
	}
	if next >= m.totalRows() {
		next = m.totalRows() - 1
	}
	m.focusRow = next

	if m.isRequiredRow(next) {
		m.focusField = kvValueField
	} else if wasRequired {
		m.focusField = kvKeyField
	}
	m.recalculate()
}

func (m *flexibleKeyValuePairListModel) moveField(delta int) {
	if m.isRequiredRow(m.focusRow) {
		m.focusField = kvValueField
		m.recalculate()
		return
	}

	next := int(m.focusField) + delta
	if next < int(kvKeyField) {
		next = int(kvKeyField)
	}
	if next > int(kvValueField) {
		next = int(kvValueField)
	}
	m.focusField = kvField(next)
	m.recalculate()
}

func (m *flexibleKeyValuePairListModel) appendText(text string) {
	row := m.focusedRow()
	if row.required || m.focusField == kvValueField {
		row.value += text
	} else {
		row.key += text
	}
	m.recalculate()
}

func (m *flexibleKeyValuePairListModel) backspace() {
	row := m.focusedRow()
	if row.required || m.focusField == kvValueField {
		row.value = trimLastRune(row.value)
	} else {
		row.key = trimLastRune(row.key)
	}
	m.recalculate()
}

func (m *flexibleKeyValuePairListModel) handleTab() {
	if !m.isRequiredRow(m.focusRow) && m.focusField == kvKeyField {
		currentIndex := m.focusRow - len(m.required)
		autocomplete := autocompleteOptionalKey(m.optional[currentIndex].key, m.schema.optional, usedOptionalKeys(m.optional, currentIndex))
		if autocomplete.advanced {
			m.optional[currentIndex].key = autocomplete.input
			m.recalculate()
			return
		}
	}

	m.advanceFocus()
}

func (m *flexibleKeyValuePairListModel) advanceFocus() {
	if m.isRequiredRow(m.focusRow) {
		if m.focusRow < len(m.required)-1 {
			m.focusRow++
			m.focusField = kvValueField
		} else {
			m.focusRow = len(m.required)
			m.focusField = kvKeyField
		}
		m.recalculate()
		return
	}

	if m.focusField == kvKeyField {
		m.focusField = kvValueField
		m.recalculate()
		return
	}

	if m.focusRow < m.totalRows()-1 {
		m.focusRow++
		if m.isRequiredRow(m.focusRow) {
			m.focusField = kvValueField
		} else {
			m.focusField = kvKeyField
		}
	}
	m.recalculate()
}

func (m *flexibleKeyValuePairListModel) recalculate() {
	m.ensureTrailingBlankRow()
	m.validateRows()
	m.updateAutocomplete()
	m.updatePreview()
}

func (m *flexibleKeyValuePairListModel) ensureTrailingBlankRow() {
	for len(m.optional) > 1 && m.isBlankOptionalRow(len(m.optional)-1) && m.isBlankOptionalRow(len(m.optional)-2) {
		m.optional = m.optional[:len(m.optional)-1]
	}

	if len(m.optional) == 0 || !m.isBlankOptionalRow(len(m.optional)-1) {
		m.optional = append(m.optional, kvRow{})
	}
}

func (m *flexibleKeyValuePairListModel) validateRows() {
	for i := range m.required {
		m.required[i].err = ""
		if strings.TrimSpace(m.required[i].value) == "" {
			m.required[i].err = "required value"
		}
	}

	seen := map[string]int{}
	for i := range m.optional {
		m.optional[i].err = ""
		key := strings.TrimSpace(m.optional[i].key)
		value := strings.TrimSpace(m.optional[i].value)

		if key == "" && value == "" {
			continue
		}
		if key == "" {
			m.optional[i].err = "missing key"
			continue
		}
		if value == "" {
			m.optional[i].err = "missing value"
		}
		if !slices.Contains(m.schema.optional, key) {
			m.optional[i].err = "unknown key"
			continue
		}
		if prev, ok := seen[key]; ok {
			m.optional[i].err = "duplicate key"
			if m.optional[prev].err == "" {
				m.optional[prev].err = "duplicate key"
			}
			continue
		}
		seen[key] = i
	}
}

func (m *flexibleKeyValuePairListModel) updateAutocomplete() {
	m.suggestion = ""
	m.candidates = nil

	if m.isRequiredRow(m.focusRow) || m.focusField != kvKeyField {
		return
	}

	currentIndex := m.focusRow - len(m.required)
	autocomplete := autocompleteOptionalKey(m.optional[currentIndex].key, m.schema.optional, usedOptionalKeys(m.optional, currentIndex))
	m.suggestion = autocomplete.suggestion
	m.candidates = autocomplete.candidates
}

func (m *flexibleKeyValuePairListModel) updatePreview() {
	payload := map[string]string{}
	incomplete := false

	for _, row := range m.required {
		if strings.TrimSpace(row.value) != "" {
			payload[row.key] = row.value
		}
		if row.err != "" {
			incomplete = true
		}
	}

	for _, row := range m.optional {
		if row.err != "" {
			if strings.TrimSpace(row.key) != "" || strings.TrimSpace(row.value) != "" {
				incomplete = true
			}
			continue
		}
		if strings.TrimSpace(row.key) != "" && strings.TrimSpace(row.value) != "" {
			payload[row.key] = row.value
		}
	}

	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		m.preview = "{\n  \"error\": \"failed to render preview\"\n}"
		m.previewNote = "preview unavailable"
		return
	}

	m.preview = string(body)
	if incomplete {
		m.previewNote = "preview omits invalid or incomplete rows"
	} else {
		m.previewNote = "all visible parameters are valid"
	}
}

func (m flexibleKeyValuePairListModel) renderRow(index int) string {
	row := m.rowAt(index)
	keySuggestion := ""
	if m.focusRow == index && m.focusField == kvKeyField && !row.required {
		keySuggestion = m.suggestion
	}

	keyField := m.renderField("key", row.key, keySuggestion, m.focusRow == index && m.focusField == kvKeyField, row.required, true)
	valueField := m.renderField("value", row.value, "", m.focusRow == index && m.focusField == kvValueField, false, false)

	statusBadge := m.renderRowStatus(row, index)

	rowLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyField,
		valueField,
	)
	if statusBadge != "" {
		rowLine = lipgloss.JoinHorizontal(lipgloss.Top, rowLine, "  ", statusBadge)
	}

	if m.focusRow == index && m.focusField == kvKeyField && len(m.candidates) > 1 {
		rowLine = lipgloss.JoinVertical(lipgloss.Left, rowLine, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")).
			Render("  suggestions: "+strings.Join(m.candidates, ", ")))
	}

	return rowLine
}

func (m flexibleKeyValuePairListModel) renderField(label string, value string, suggestion string, focused bool, locked bool, isKey bool) string {
	labelStyle := lipgloss.NewStyle()
	valueStyle := lipgloss.NewStyle()
	ghostStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00eeff")).Background(lipgloss.Color("#b0b0b0"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#F8FAFC"))
	boxStyle := lipgloss.NewStyle().
		Padding(0, 1)
	fieldWidth := kvValueFieldWidth
	contentStyle := lipgloss.NewStyle()

	if isKey {
		boxStyle = boxStyle.Background(lipgloss.Color("#b0b0b0"))
		contentStyle = contentStyle.Background(lipgloss.Color("#b0b0b0"))
		fieldWidth = kvKeyFieldWidth
		labelStyle = labelStyle.Foreground(lipgloss.Color("#b5d8ff")).Background(lipgloss.Color("#b0b0b0"))
		valueStyle = valueStyle.Foreground(lipgloss.Color("#DBEAFE"))
	} else {
		boxStyle = boxStyle.Background(lipgloss.Color("#7d7d7d"))
		contentStyle = contentStyle.Background(lipgloss.Color("#7d7d7d"))
		labelStyle = labelStyle.Foreground(lipgloss.Color("#b5d8ff")).Background(lipgloss.Color("#b0b0b0"))
		valueStyle = valueStyle.Foreground(lipgloss.Color("#FEF3C7"))
	}

	if locked {
		valueStyle = valueStyle.Bold(true)
	}
	if focused {
		boxStyle = boxStyle.
			BorderForeground(lipgloss.Color("#F8FAFC"))
		valueStyle = valueStyle.Bold(true)
	}

	displayValue := ""
	if value == "" {
		displayValue = lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).Render(label)
	} else {
		displayValue = valueStyle.Render(value)
	}

	cursor := ""
	if focused {
		cursor = cursorStyle.Render(" ")
	}

	baseContent := lipgloss.JoinHorizontal(lipgloss.Top, displayValue, cursor, ghostStyle.Render(suggestion))
	content := contentStyle.
		Width(fieldWidth).
		MaxWidth(fieldWidth).
		Render(baseContent)

	return boxStyle.
		Width(fieldWidth).
		MaxWidth(fieldWidth).
		Render(content)
}

func (m flexibleKeyValuePairListModel) renderRowStatus(row kvRow, index int) string {
	statuses := make([]string, 0, 2)

	if row.required {
		label := "required"
		color := lipgloss.Color("#93C5FD")
		if row.err == "required value" {
			label = "required value"
			color = lipgloss.Color("#FCA5A5")
		}
		statuses = append(statuses, lipgloss.NewStyle().
			Foreground(color).
			Render(label))
	} else if index == m.totalRows()-1 && strings.TrimSpace(row.key) == "" && strings.TrimSpace(row.value) == "" {
		statuses = append(statuses, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#93C5FD")).
			Render("optional"))
	}

	if row.err != "" && (!row.required || row.err != "required value") {
		statuses = append(statuses, m.renderValidationBadge(row.err))
	}

	if len(statuses) == 0 {
		return ""
	}

	return strings.Join(statuses, "  ")
}

func (m flexibleKeyValuePairListModel) renderValidationBadge(err string) string {
	if err == "unknown key" {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B91C1C")).
			Background(lipgloss.Color("#FEE2E2")).
			Padding(0, 1).
			Render(err)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FCA5A5")).
		Render(err)
}

func (m flexibleKeyValuePairListModel) renderPreview() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#86EFAC")).
		Render("Preview")

	note := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")).
		Render(m.previewNote)

	body := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#14532D")).
		Padding(0, 1).
		Foreground(lipgloss.Color("#DCFCE7")).
		Render(m.preview)

	return strings.Join([]string{title, note, body}, "\n")
}

func (m flexibleKeyValuePairListModel) totalRows() int {
	return len(m.required) + len(m.optional)
}

func (m flexibleKeyValuePairListModel) isRequiredRow(index int) bool {
	return index < len(m.required)
}

func (m flexibleKeyValuePairListModel) rowAt(index int) kvRow {
	if m.isRequiredRow(index) {
		return m.required[index]
	}
	return m.optional[index-len(m.required)]
}

func (m *flexibleKeyValuePairListModel) focusedRow() *kvRow {
	if m.isRequiredRow(m.focusRow) {
		return &m.required[m.focusRow]
	}
	return &m.optional[m.focusRow-len(m.required)]
}

func (m flexibleKeyValuePairListModel) isBlankOptionalRow(index int) bool {
	row := m.optional[index]
	return strings.TrimSpace(row.key) == "" && strings.TrimSpace(row.value) == ""
}

func kvDefaultSchema() kvSchema {
	return kvSchema{
		required: []string{"endpoint", "method"},
		optional: []string{"query", "header", "timeout", "retries", "profile", "region", "token", "output"},
	}
}

func autocompleteOptionalKey(input string, allowed []string, used []string) kvAutocomplete {
	candidates := make([]string, 0)
	usedSet := make(map[string]struct{}, len(used))
	for _, key := range used {
		usedSet[key] = struct{}{}
	}

	for _, key := range allowed {
		if _, taken := usedSet[key]; taken {
			continue
		}
		if strings.HasPrefix(key, input) {
			candidates = append(candidates, key)
		}
	}
	slices.Sort(candidates)

	if len(candidates) == 0 {
		return kvAutocomplete{input: input}
	}
	if len(candidates) == 1 {
		suggestion := ""
		if strings.HasPrefix(candidates[0], input) {
			suggestion = candidates[0][len(input):]
		}
		return kvAutocomplete{
			input:      candidates[0],
			suggestion: suggestion,
			candidates: candidates,
			advanced:   candidates[0] != input,
		}
	}

	common := longestCommonPrefix(candidates)
	suggestion := ""
	if strings.HasPrefix(common, input) {
		suggestion = common[len(input):]
	}

	return kvAutocomplete{
		input:      commonIfLonger(common, input),
		suggestion: suggestion,
		candidates: candidates,
		advanced:   len(common) > len(input),
	}
}

func usedOptionalKeys(rows []kvRow, skip int) []string {
	out := make([]string, 0)
	for i, row := range rows {
		if i == skip {
			continue
		}
		key := strings.TrimSpace(row.key)
		if key != "" {
			out = append(out, key)
		}
	}
	return out
}

func commonIfLonger(candidate string, current string) string {
	if len(candidate) > len(current) {
		return candidate
	}
	return current
}

func trimLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	return string(runes[:len(runes)-1])
}

func longestCommonPrefix(values []string) string {
	if len(values) == 0 {
		return ""
	}

	prefix := values[0]
	for _, value := range values[1:] {
		for !strings.HasPrefix(value, prefix) && prefix != "" {
			prefix = prefix[:len(prefix)-1]
		}
		if prefix == "" {
			return ""
		}
	}
	return prefix
}
