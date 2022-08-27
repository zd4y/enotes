package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	bold      = lipgloss.NewStyle().Bold(true)
	keyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))
	sepStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	dot       = sepStyle.Render(" • ")
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
		newHeaderHeight := lipgloss.Height(m.noteHeaderView())
		m.noteViewport.Height -= newHeaderHeight - 1
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

	return fmt.Sprintf("%s\n%s\n%s", m.noteHeaderView(), m.noteViewport.View(), m.noteFooterView())
}

func (m model) noteHeaderView() string {
	item, ok := m.list.SelectedItem().(fileItem)
	if !ok {
		return ""
	}
	title := fmt.Sprintf("File: %s", bold.Render(item.file.Name()))
	return title
}

func (m model) noteFooterView() string {
	return strings.Join([]string{
		help("↑/k", "up"),
		help("↓/j", "down"),
		help("e", "edit note"),
		help("q/esc", "go back"),
		help("ctrl+c", "quit"),
	}, dot)
}

func help(key, desc string) string {
	return keyStyle.Render(key) + " " + descStyle.Render(desc)
}
