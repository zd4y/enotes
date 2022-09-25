package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zd4y/enotes/enotes"
)

func newNoteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			m.noteAlreadyExists = false
			m.resetChosen()
			return m, nil
		case "enter":
			m.noteAlreadyExists = false
			newNoteName := m.textInput.Value()
			if newNoteName == "" {
				newNoteName = time.Now().Format(time.Stamp)
			}
			if ok, err := enotes.NoteExists(newNoteName); ok {
				m.noteAlreadyExists = true
				return m, nil
			} else if err != nil {
				m.err = err
				return m, nil
			}
			m.newNoteName = newNoteName
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func newNoteView(m model) string {
	s := fmt.Sprintf(
		"New note name?\n\n%s\n",
		m.textInput.View(),
	)

	if m.noteAlreadyExists {
		s += "\nNote already exists"
	}

	return s + "\n(esc to quit)\n"
}
