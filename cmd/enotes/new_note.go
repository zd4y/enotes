package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zd4y/enotes/pkg/enotes"
)

func newNoteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			m.resetChosen()
			return m, nil
		case "enter":
			newNotePath := m.textInput.Value()
			if newNotePath == "" {
				newNotePath = time.Now().Format(time.Stamp) + ".md.age"
			}
			if ok, err := enotes.NoteExists(newNotePath); ok {
				return m, tea.Println("Note already exists")
			} else if err != nil {
				m.err = err
				return m, nil
			}
			m.newNotePath = newNotePath
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func newNoteView(m model) string {
	return fmt.Sprintf(
		"New note filename?\n\n%s\n\n%s\n",
		m.textInput.View(),
		"(esc to quit)",
	)
}
