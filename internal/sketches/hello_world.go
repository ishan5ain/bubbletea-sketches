package sketches

import tea "github.com/charmbracelet/bubbletea"

type helloWorldModel struct{}

func NewHelloWorld() tea.Model {
	return helloWorldModel{}
}

func (m helloWorldModel) Init() tea.Cmd {
	return nil
}

func (m helloWorldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

func (m helloWorldModel) View() string {
	return "Hello, World!\n\nPress q to quit.\n"
}
