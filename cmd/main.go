package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/zd4y/notes/pkg/notes"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type fileItem struct {
	file fs.FileInfo
}

func (i fileItem) Title() string       { return i.file.Name() }
func (i fileItem) Description() string { return i.file.ModTime().String() }
func (i fileItem) FilterValue() string { return i.file.Name() }

type model struct {
	list list.Model
	chosen int
	editorActive bool
	newNoteName string
	password string
	passwordCorrect bool
	passwordVerified bool
	textInput textinput.Model
	err error
}

func initialModel() model {
	items := []list.Item{
		item{title: "New note", desc: "Write a new encrypted note"},
	}

	textInput := textinput.New()
	textInput.Placeholder = "Password"
	textInput.Focus()
	m := model{chosen: -1, list: list.New(items, list.NewDefaultDelegate(), 0, 0), textInput: textInput}
	m.list.Title = "Notes"
	return m
}

func (m model) inPassword() bool {
	return len(m.password) == 0
}

func (m model) inNote() bool {
	return m.chosen > 0
}

func (m model) inNewNote() bool {
	return m.chosen == 0
}

func (m model) inNewNoteEditor() bool {
	return len(m.newNoteName) > 0
}

func (m *model) resetChosen() {
	m.chosen = -1
}

func (m *model) toNewNote() {
	m.chosen = 0
}

func (m *model) toNote(inIndex int) {
	m.chosen = inIndex
}

type editorFinishedMsg struct{ err error }

func openEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

type verifyPasswordMsg struct {
	passwordsMatch bool
	err error
}

func verifyPassword(password string) tea.Cmd {
	return func() tea.Msg {
		err := notes.VerifyPassword(password)
		return verifyPasswordMsg { err == nil, err }
	}
}

type dirFilesMsg struct { files []fs.FileInfo }

func getDirFiles() tea.Msg {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	return dirFilesMsg{files}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, getDirFiles)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.editorActive = false
	case dirFilesMsg:
		index := len(m.list.Items())
		var cmd tea.Cmd
		for i, file := range msg.files {
			cmd = m.list.InsertItem(index + i, fileItem{file})
		}
		return m, cmd
	case verifyPasswordMsg:
		m.passwordCorrect = msg.passwordsMatch
		m.passwordVerified = true
		m.err = msg.err
		return m, nil
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.inPassword() {
		return passwordUpdate(msg, m)
	}
	if !m.passwordCorrect {
		return m, nil
	}
	if m.inNote() {
		return noteUpdate(msg, m)
	}
	if m.inNewNoteEditor() {
		return newNoteEditorUpdate(msg, m)
	}
	if m.inNewNote() {
		return newNoteUpdate(msg, m)
	}
	return fileListUpdate(msg, m)
}

func (m model) View() string {
	if m.inPassword() {
		return passwordView(m)
	}

	if !m.passwordCorrect {
		if m.passwordVerified {
			return "incorrect password or failed verifying password:\n\t" + m.err.Error()
		} else {
			return "verifying password..."
		}
	}

	if m.inNote() {
		return noteView(m)
	}
	if m.inNewNoteEditor() {
		return newNoteEditorView(m)
	}
	if m.inNewNote() {
		return newNoteView(m)
	}
	return fileListView(m)
}

func passwordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			return m, tea.Quit
    case "enter":
			m.password = m.textInput.Value()
      return m, verifyPassword(m.password)
    }
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func noteUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc", "q":
			m.resetChosen()
			return m, nil
		case "e":
			return m, openEditor()
    }
	}
	return m, nil
}

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

func fileListUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg := msg.String(); msg {
		case "enter":
			index := m.list.Index()
			if index == 0 {
				m.textInput = textinput.New()
				m.textInput.Placeholder = "File name (leave empty for current date)"
				m.textInput.Focus()
				m.toNewNote()
			} else {
				m.toNote(index)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func passwordView(m model) string {
	return fmt.Sprintf(
		"Password?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func noteView(m model) string {
	listItem := m.list.Items()[m.chosen]
	switch item := listItem.(type) {
	case item:
		return fmt.Sprintf("in note %s", item.title)
	case fileItem:
		return fmt.Sprintf("in note %s", item.file.Name())
	}

	panic("unhandled item type")
}

func newNoteEditorView(m model) string {
	return "loading editor"
}

func newNoteView(m model) string {
	return fmt.Sprintf(
		"New note filename?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func fileListView(m model) string {
	return docStyle.Render(m.list.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
