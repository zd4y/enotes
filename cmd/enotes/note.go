package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func noteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openNoteMsg:
		m.loadingNote = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.noteContents = msg.note
	case tea.KeyMsg:
		switch msg := msg.String(); msg {
		case "esc", "q":
			m.resetChosen()
			return m, nil
		case "e":
			if !m.loadingNote {
				m.editorActive = true
				item := m.list.SelectedItem().(fileItem)
				return m, editNote(item.file.Name(), m.password)
			}
		}
	}

	r, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(m.width))
	if err != nil {
		m.err = err
	}
	out, err := r.Render(m.noteContents)
	if err != nil {
		m.err = err
	}
	m.noteViewport.SetContent(out)

	var cmd tea.Cmd
	m.noteViewport, cmd = m.noteViewport.Update(msg)
	return m, cmd
}

func noteView(m model) string {
	if m.loadingNote {
		return fmt.Sprintf("%s Decrypting note\n", m.spinner.View())
	}

	if m.editorActive {
		return fmt.Sprintf("%s Loading editor\n", m.spinner.View())
	}

	if m.err != nil {
		return "error loading note: " + m.err.Error() + "\n"
	}

	return m.noteViewport.View()
}
