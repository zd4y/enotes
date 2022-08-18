package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"

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
	textInput textinput.Model
	err error
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
		for i, file := range msg.files {
			m.list.InsertItem(index + i, fileItem{file})
		}
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
			v := m.textInput.Value()
			m.newNoteName = v
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
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
	items := []list.Item{
		item{title: "New note", desc: "Write a new encrypted note"},
	}

	m := model{chosen: -1, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
