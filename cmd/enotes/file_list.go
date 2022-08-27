package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func fileListUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg := msg.String(); msg {
		case "q", "esc":
			m.quitting = true
			return m, nil
		case "enter":
			index := m.list.Index()
			if index == 0 {
				m.textInput = textinput.New()
				m.textInput.Placeholder = "New note name (leave empty for current date)"
				m.textInput.Focus()
				m.toNewNote()
				return m, textinput.Blink
			} else {
				m.toNote(index)
				item := m.list.SelectedItem().(fileItem)
				m.loadingNote = true
				return m, openNote(item.file.Name(), m.password)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func fileListView(m model) string {
	return docStyle.Render(m.list.View())
}
