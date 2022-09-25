package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func newNoteEditorUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if m.editorActive {
		return m, nil
	}
	m.editorActive = true
	return m, createNote(m.newNoteName, m.password)
}

func newNoteEditorView(m model) string {
	return fmt.Sprintf("%s Loading editor\n", m.spinner.View())
}
