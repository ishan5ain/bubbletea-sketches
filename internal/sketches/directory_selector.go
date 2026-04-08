package sketches

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type directoryCompletion struct {
	Input       string
	Suggestion  string
	Matches     []string
	Error       string
	ResolvedAbs string
}

type directorySelectorModel struct {
	input        string
	cwd          string
	home         string
	suggestion   string
	matches      []string
	selectedPath string
	errorMsg     string
}

func NewDirectorySelector() tea.Model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = cwd
	}

	model := directorySelectorModel{
		cwd:  cwd,
		home: home,
	}
	model.refreshCompletion()
	return model
}

func (m directorySelectorModel) Init() tea.Cmd {
	return nil
}

func (m directorySelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = trimLastRune(m.input)
				m.errorMsg = ""
				m.selectedPath = ""
				m.refreshCompletion()
			}
		case "tab":
			m.applyCompletion()
		case "enter":
			m.confirmSelection()
		default:
			if text := msg.Key().Text; text != "" {
				m.input += text
				m.errorMsg = ""
				m.selectedPath = ""
				m.refreshCompletion()
			}
		}
	}

	return m, nil
}

func (m directorySelectorModel) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Background(lipgloss.Color("#0F172A")).
		Padding(0, 0, 1, 0).
		Render("directory-selector")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")).
		Render("type a path, tab to autocomplete directories, enter to confirm")

	prompt := renderPromptLine(m.input, m.suggestion)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Render("\ntab autocomplete | backspace delete | enter select | q quit")

	lines := []string{title, subtitle, prompt}

	if len(m.matches) > 1 {
		lines = append(lines, renderMatchList(m.matches))
	}

	if m.selectedPath != "" {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#86EFAC")).
			Render("Selected: "+m.selectedPath))
	}

	if m.errorMsg != "" {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FCA5A5")).
			Render(m.errorMsg))
	}

	lines = append(lines, help)

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(1, 2).
		Render(strings.Join(lines, "\n"))

	return tea.NewView(panel)
}

func (m *directorySelectorModel) applyCompletion() {
	completion := completeDirectoryInput(m.input, m.cwd, m.home)
	m.input = completion.Input
	m.suggestion = completion.Suggestion
	m.matches = completion.Matches
	m.errorMsg = completion.Error
	m.selectedPath = ""
}

func (m *directorySelectorModel) confirmSelection() {
	absPath, err := resolveDirectoryInput(m.input, m.cwd, m.home)
	if err != nil {
		m.selectedPath = ""
		m.errorMsg = err.Error()
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		m.selectedPath = ""
		m.errorMsg = fmt.Sprintf("cannot access %q: %v", absPath, err)
		return
	}
	if !info.IsDir() {
		m.selectedPath = ""
		m.errorMsg = fmt.Sprintf("%q is not a directory", absPath)
		return
	}

	m.selectedPath = absPath
	m.errorMsg = ""
	m.refreshCompletion()
}

func (m *directorySelectorModel) refreshCompletion() {
	completion := analyzeDirectoryInput(m.input, m.cwd, m.home)
	m.suggestion = completion.Suggestion
	m.matches = completion.Matches
	if m.errorMsg == "" {
		m.errorMsg = completion.Error
	}
}

func renderPromptLine(input string, suggestion string) string {
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#38BDF8")).
		Bold(true)

	typedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8FAFC"))

	ghostStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B"))

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0F172A")).
		Background(lipgloss.Color("#F8FAFC"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		promptStyle.Render("cd "),
		typedStyle.Render(input),
		ghostStyle.Render(suggestion),
		cursorStyle.Render(" "),
	)
}

func renderMatchList(matches []string) string {
	items := make([]string, 0, len(matches))
	for _, match := range matches {
		items = append(items, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Render(match))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")).
		Render("Matches: ")

	return label + strings.Join(items, "  ")
}

func completeDirectoryInput(input string, cwd string, home string) directoryCompletion {
	completion := analyzeDirectoryInput(input, cwd, home)
	if completion.Error != "" || len(completion.Matches) == 0 {
		return completion
	}

	if len(completion.Matches) == 1 {
		completion.Input = ensureTrailingSeparator(completion.Matches[0])
		completion = analyzeDirectoryInput(completion.Input, cwd, home)
		return completion
	}

	common := longestCommonPrefix(stripTrailingSeparators(completion.Matches))
	if len(common) > len(input) {
		completion.Input = common
	}
	completion = analyzeDirectoryInput(completion.Input, cwd, home)
	return completion
}

func analyzeDirectoryInput(input string, cwd string, home string) directoryCompletion {
	searchDir, fragment, displayPrefix, err := directoryCompletionContext(input, cwd, home)
	if err != nil {
		return directoryCompletion{
			Input: input,
			Error: err.Error(),
		}
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return directoryCompletion{
			Input: input,
			Error: fmt.Sprintf("cannot read %q: %v", searchDir, err),
		}
	}

	matches := make([]string, 0)
	names := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, fragment) {
			names = append(names, name)
			matches = append(matches, displayPrefix+name)
		}
	}

	slices.Sort(matches)
	slices.Sort(names)

	suggestion := ""
	if len(names) == 1 {
		full := ensureTrailingSeparator(displayPrefix + names[0])
		if strings.HasPrefix(full, input) {
			suggestion = full[len(input):]
		}
	} else if len(names) > 1 {
		common := displayPrefix + longestCommonPrefix(names)
		if strings.HasPrefix(common, input) {
			suggestion = common[len(input):]
		}
	}

	return directoryCompletion{
		Input:       input,
		Suggestion:  suggestion,
		Matches:     appendDirectorySeparators(matches),
		ResolvedAbs: resolvePathRaw(input, cwd, home),
	}
}

func directoryCompletionContext(input string, cwd string, home string) (searchDir string, fragment string, displayPrefix string, err error) {
	if input == "~" {
		return filepath.Dir(home), filepath.Base(home), "", nil
	}

	resolved := resolvePathRaw(input, cwd, home)
	if strings.HasSuffix(input, string(filepath.Separator)) {
		return resolved, "", input, nil
	}

	searchDir = filepath.Dir(resolved)
	fragment = filepath.Base(resolved)
	if input == "" {
		searchDir = cwd
		fragment = ""
	}

	if idx := strings.LastIndex(input, string(filepath.Separator)); idx >= 0 {
		displayPrefix = input[:idx+1]
	}

	return searchDir, fragment, displayPrefix, nil
}

func resolveDirectoryInput(input string, cwd string, home string) (string, error) {
	absPath := resolvePathRaw(input, cwd, home)
	absPath = filepath.Clean(absPath)
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(cwd, absPath)
	}
	return absPath, nil
}

func resolvePathRaw(input string, cwd string, home string) string {
	switch {
	case input == "":
		return cwd
	case input == "~":
		return home
	case strings.HasPrefix(input, "~"+string(filepath.Separator)):
		return filepath.Join(home, input[2:])
	case filepath.IsAbs(input):
		return input
	default:
		return filepath.Join(cwd, input)
	}
}

func appendDirectorySeparators(matches []string) []string {
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		out = append(out, ensureTrailingSeparator(match))
	}
	return out
}

func stripTrailingSeparators(matches []string) []string {
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		out = append(out, strings.TrimSuffix(match, string(filepath.Separator)))
	}
	return out
}

func ensureTrailingSeparator(value string) string {
	if strings.HasSuffix(value, string(filepath.Separator)) {
		return value
	}
	return value + string(filepath.Separator)
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
