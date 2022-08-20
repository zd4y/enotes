package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func newNoteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			m.resetChosen()
			return m, nil
    case "enter":
			m.newNoteName = m.textInput.Value()
      return m, nil
    }
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func newNoteView(m model) string {
	return fmt.Sprintf(
		"New note filename?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

