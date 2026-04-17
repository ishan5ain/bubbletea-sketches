package sketchbookapp

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/ishansain/bubbletea-sketches/internal/sketches"
)

const pageSize = 10

type sketchItem struct {
	title      string
	sketchName string
}

func (i sketchItem) Title() string {
	return i.title
}

func (i sketchItem) Description() string {
	return ""
}

func (i sketchItem) FilterValue() string {
	return i.sketchName
}

type browserModel struct {
	names     []string
	page      int
	pageCount int
	width     int
	height    int
	list      list.Model
	selected  string
}

func RunBrowser(stdin io.Reader, stdout io.Writer) error {
	finalModel, err := runProgram(newBrowserModel(sketches.Names()), stdin, stdout)
	if err != nil {
		return err
	}

	browser, ok := finalModel.(browserModel)
	if !ok || browser.selected == "" {
		return nil
	}

	factory, ok := sketches.Get(browser.selected)
	if !ok {
		return fmt.Errorf("unknown sketch %q", browser.selected)
	}

	_, err = runProgram(factory(), stdin, stdout)
	return err
}

func RunLegacy(args []string, stdin io.Reader, stdout io.Writer) error {
	model, err := SelectSketch(args)
	if err != nil {
		return err
	}

	_, err = runProgram(model, stdin, stdout)
	return err
}

func SelectSketch(args []string) (tea.Model, error) {
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

func newBrowserModel(names []string) browserModel {
	model := browserModel{
		names:     append([]string(nil), names...),
		pageCount: pageCount(len(names)),
	}
	model.setPage(0)
	return model
}

func (m browserModel) Init() tea.Cmd {
	return nil
}

func (m browserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeList()
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if selected, ok := m.list.SelectedItem().(sketchItem); ok {
				m.selected = selected.sketchName
				return m, tea.Quit
			}
			return m, nil
		case "left":
			m.setPage(m.page - 1)
			return m, nil
		case "right":
			m.setPage(m.page + 1)
			return m, nil
		}

		if idx, ok := digitIndex(msg.String()); ok {
			if sketchName, ok := m.pageSketchName(idx); ok {
				m.selected = sketchName
				return m, tea.Quit
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m browserModel) View() tea.View {
	headerText := "Sketchbook"
	if m.pageCount > 0 {
		headerText = fmt.Sprintf("Sketchbook [%d/%d]", m.page+1, m.pageCount)
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Render(headerText)

	body := m.list.View()
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Render("0-9 run | left/right page | q quit")

	return tea.NewView(strings.Join([]string{header, body, footer}, "\n"))
}

func (m *browserModel) setPage(page int) {
	if m.pageCount == 0 {
		m.page = 0
		m.list = list.New(nil, list.NewDefaultDelegate(), m.width, bodyHeight(m.height))
		m.list.SetShowTitle(false)
		m.list.SetShowStatusBar(false)
		m.list.SetFilteringEnabled(false)
		m.list.SetShowHelp(false)
		m.list.SetShowPagination(false)
		m.list.Title = "Sketchbook"
		return
	}

	if page < 0 {
		page = 0
	}
	if page >= m.pageCount {
		page = m.pageCount - 1
	}

	m.page = page
	m.list = newPageList(m.names, m.page, m.width, m.height)
	m.list.Title = fmt.Sprintf("Sketchbook [%d/%d]", m.page+1, m.pageCount)
}

func (m *browserModel) resizeList() {
	m.setPage(m.page)
}

func (m browserModel) pageSketchName(index int) (string, bool) {
	start := m.page * pageSize
	pos := start + index
	if pos < 0 || pos >= len(m.names) {
		return "", false
	}
	if index < 0 || index >= pageSize {
		return "", false
	}
	return m.names[pos], true
}

func pageCount(total int) int {
	if total == 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func pageItems(names []string, page int) []list.Item {
	start := page * pageSize
	if start >= len(names) {
		return nil
	}

	end := start + pageSize
	if end > len(names) {
		end = len(names)
	}

	items := make([]list.Item, 0, end-start)
	for i, name := range names[start:end] {
		items = append(items, sketchItem{
			title:      fmt.Sprintf("%d %s", i, name),
			sketchName: name,
		})
	}

	return items
}

func newPageList(names []string, page int, width int, height int) list.Model {
	items := pageItems(names, page)
	l := list.New(items, list.NewDefaultDelegate(), width, bodyHeight(height))
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	return l
}

func bodyHeight(height int) int {
	body := height - 2
	if body < 1 {
		return 1
	}
	return body
}

func digitIndex(value string) (int, bool) {
	if len(value) != 1 {
		return 0, false
	}
	ch := value[0]
	if ch < '0' || ch > '9' {
		return 0, false
	}
	return int(ch - '0'), true
}

func runProgram(model tea.Model, stdin io.Reader, stdout io.Writer) (tea.Model, error) {
	program := tea.NewProgram(model, tea.WithInput(stdin), tea.WithOutput(stdout))
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}
	return finalModel, nil
}
