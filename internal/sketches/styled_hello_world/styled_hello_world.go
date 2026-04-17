package styled_hello_world

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

type styleSection struct {
	name        string
	description string
	render      func() string
}

type styledHelloWorldModel struct {
	index    int
	sections []styleSection
}

func NewStyledHelloWorld() tea.Model {
	return styledHelloWorldModel{
		sections: styledHelloWorldSections(),
	}
}

func (m styledHelloWorldModel) Init() tea.Cmd {
	return nil
}

func (m styledHelloWorldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			m.index--
			if m.index < 0 {
				m.index = len(m.sections) - 1
			}
		case "right", "l":
			m.index++
			if m.index >= len(m.sections) {
				m.index = 0
			}
		}
	}

	return m, nil
}

func (m styledHelloWorldModel) View() tea.View {
	section := m.sections[m.index]

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#1D4ED8")).
		Padding(0, 1).
		Render("styled-hello-world")

	meta := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#475569")).
		Render(fmt.Sprintf("Section %d/%d", m.index+1, len(m.sections)))

	heading := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Render(section.name)

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#566e90")).
		Render(section.description)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#df8918")).
		Render("left/h previous  right/l next  q quit")

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2563EB")).
		Padding(1, 2).
		MarginTop(1).
		Render(section.render())

	content := strings.Join([]string{
		lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", meta),
		heading,
		description,
		panel,
		help,
	}, "\n")

	return tea.NewView(content)
}

func styledHelloWorldSections() []styleSection {
	return []styleSection{
		{
			name:        "Text Emphasis",
			description: "Bold, italic, faint, blink, reverse, and strikethrough applied to the same greeting.",
			render: func() string {
				samples := []string{
					labeledSample("Bold", lipgloss.NewStyle().Bold(true)),
					labeledSample("Italic", lipgloss.NewStyle().Italic(true)),
					labeledSample("Faint", lipgloss.NewStyle().Faint(true)),
					labeledSample("Blink", lipgloss.NewStyle().Blink(true)),
					labeledSample("Reverse", lipgloss.NewStyle().Reverse(true)),
					labeledSample("Strikethrough", lipgloss.NewStyle().Strikethrough(true)),
				}
				return strings.Join(samples, "\n")
			},
		},
		{
			name:        "Colors",
			description: "ANSI, 256-color, truecolor, and adaptive colors show how Lip Gloss picks terminal-friendly palettes.",
			render: func() string {
				samples := []string{
					labeledSample("ANSI 16", lipgloss.NewStyle().Foreground(lipgloss.Color("9"))),
					labeledSample("ANSI 256", lipgloss.NewStyle().Foreground(lipgloss.Color("201"))),
					labeledSample("Truecolor", lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Background(lipgloss.Color("#082F49"))),
					labeledSample("Adaptive", lipgloss.NewStyle().Foreground(compat.AdaptiveColor{
						Light: lipgloss.Color("#1D4ED8"),
						Dark:  lipgloss.Color("#93C5FD"),
					})),
				}
				return strings.Join(samples, "\n")
			},
		},
		{
			name:        "Underline Decoration",
			description: "Underline styles and underline colors vary by terminal, but Lip Gloss v2 exposes the main decoration controls directly.",
			render: func() string {
				samples := []string{
					labeledSample("Single", lipgloss.NewStyle().Underline(true)),
					labeledSample("Double", lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDouble)),
					labeledSample("Curly", lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineCurly)),
					labeledSample("Dotted", lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDotted)),
					labeledSample("Dashed + Color", lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDashed).UnderlineColor(lipgloss.Color("#F97316"))),
					labeledSample("Underline Spaces", lipgloss.NewStyle().Underline(true).UnderlineSpaces(true).Width(20)),
				}
				return strings.Join(samples, "\n")
			},
		},
		{
			name:        "Spacing",
			description: "Padding and margin shape the breathing room around the greeting without changing the underlying text.",
			render: func() string {
				card := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#F8FAFC")).
					Background(lipgloss.Color("#0F766E")).
					Padding(1, 3).
					Margin(1, 2, 1, 4)
				compact := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#111827")).
					Background(lipgloss.Color("#FDE68A")).
					Padding(0, 1)
				return strings.Join([]string{
					"Margin adds outer space, padding adds inner space:",
					card.Render("Hello, World!"),
					compact.Render("Hello, World!"),
				}, "\n")
			},
		},
		{
			name:        "Sizing and Alignment",
			description: "Width, height, and alignment let the same content occupy space intentionally within a frame.",
			render: func() string {
				left := lipgloss.NewStyle().
					Width(22).
					Height(5).
					Align(lipgloss.Left).
					AlignVertical(lipgloss.Top).
					Border(lipgloss.NormalBorder()).
					Padding(0, 1)
				center := lipgloss.NewStyle().
					Width(22).
					Height(5).
					Align(lipgloss.Center).
					AlignVertical(lipgloss.Center).
					Border(lipgloss.NormalBorder()).
					Padding(0, 1)
				right := lipgloss.NewStyle().
					Width(22).
					Height(5).
					Align(lipgloss.Right).
					AlignVertical(lipgloss.Bottom).
					Border(lipgloss.NormalBorder()).
					Padding(0, 1)
				return lipgloss.JoinHorizontal(
					lipgloss.Top,
					labeledBlock("Left / Top", left.Render("Hello, World!")),
					labeledBlock("Center / Middle", center.Render("Hello, World!")),
					labeledBlock("Right / Bottom", right.Render("Hello, World!")),
				)
			},
		},
		{
			name:        "Borders",
			description: "Border presets plus custom border colors define the container around the greeting.",
			render: func() string {
				rounded := borderSample("Rounded", lipgloss.RoundedBorder(), lipgloss.Color("#7C3AED"), nil)
				double := borderSample("Double", lipgloss.DoubleBorder(), lipgloss.Color("#DC2626"), lipgloss.Color("#FEE2E2"))
				thick := borderSample("Thick", lipgloss.ThickBorder(), lipgloss.Color("#047857"), nil)
				return lipgloss.JoinHorizontal(lipgloss.Top, rounded, double, thick)
			},
		},
		{
			name:        "Composition",
			description: "Inherited styles carry visual traits forward, while unset operations strip them back for targeted overrides.",
			render: func() string {
				base := lipgloss.NewStyle().
					Bold(true).
					Italic(true).
					Foreground(lipgloss.Color("#F8FAFC")).
					Background(lipgloss.Color("#7C3AED")).
					Padding(0, 1)
				inherited := lipgloss.NewStyle().
					Inherit(base).
					Underline(true)
				reset := inherited.
					UnsetItalic().
					UnsetBackground().
					Foreground(lipgloss.Color("#7C3AED"))
				return strings.Join([]string{
					labeledSample("Base", base),
					labeledSample("Inherited + Underline", inherited),
					labeledSample("Unset Italic/Background", reset),
				}, "\n")
			},
		},
	}
}

func labeledSample(label string, style lipgloss.Style) string {
	labelStyle := lipgloss.NewStyle().
		Width(26).
		Foreground(lipgloss.Color("#64748B"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(label),
		style.Render("Hello, World!"),
	)
}

func labeledBlock(label string, content string) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#334155")).
		Render(label)

	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

func borderSample(label string, border lipgloss.Border, borderColor color.Color, background color.Color) string {
	style := lipgloss.NewStyle().
		Width(18).
		Height(5).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Border(border).
		BorderForeground(borderColor).
		Padding(0, 1)
	if background != nil {
		style = style.Background(background)
	}

	return labeledBlock(label, style.Render("Hello, World!"))
}
