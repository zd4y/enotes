package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func noteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openNoteMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.loadingNote = false
		m.noteContents = msg.note

		out, err := glamour.Render(m.noteContents, "dark")
		if err != nil {
			m.err = err
		}
		m.noteViewport.SetContent(out)

		return m, nil
	case tea.KeyMsg:
		switch msg := msg.String(); msg {
		case "esc", "q":
			m.resetChosen()
			return m, nil
		case "e":
			return m, openEditor()
		}
	}

	var cmd tea.Cmd
	m.noteViewport, cmd = m.noteViewport.Update(msg)
	return m, cmd
}

func noteView(m model) string {
	if m.loadingNote {
		return fmt.Sprintf("%s Decrypting note\n", m.spinner.View())
	}

	if len(m.noteContents) > 0 {
		return m.noteViewport.View()
	}

	if m.err != nil {
		return "error loading note: " + m.err.Error()
	}

	panic("called note view without loading, contents or error")
}
