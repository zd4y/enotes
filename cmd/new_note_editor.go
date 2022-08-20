package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func newNoteEditorUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if _, ok := msg.(editorFinishedMsg); ok {
		m.resetChosen()
		m.newNoteName = ""
		return m, nil
	}
	if m.editorActive {
		return m, nil
	}
	m.editorActive = true
	return m, openEditor()
}

func newNoteEditorView(m model) string {
	return "loading editor"
}
