package sketches

import tea "charm.land/bubbletea/v2"

type helloWorldModel struct{}

func NewHelloWorld() tea.Model {
	return helloWorldModel{}
}

func (m helloWorldModel) Init() tea.Cmd {
	return nil
}

func (m helloWorldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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

func (m helloWorldModel) View() tea.View {
	return tea.NewView("Hello, World!\n\nPress q to quit.\n")
}
