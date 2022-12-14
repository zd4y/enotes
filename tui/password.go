package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func passwordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			m.quitting = true
			return m, nil
		case "enter":
			m.password = m.textInput.Value()
			return m, verifyPassword(m.password)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func passwordView(m model) string {
	return fmt.Sprintf(
		"Password?\n\n%s\n\n%s\n",
		m.textInput.View(),
		"(esc to quit)",
	)
}
