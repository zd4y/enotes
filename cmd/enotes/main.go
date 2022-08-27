package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zd4y/enotes/pkg/enotes"
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
	quitting            bool
	width               int
	list                list.Model
	chosen              int
	editorActive        bool
	newNotePath         string
	password            string
	passwordVerified    bool
	noteContents        string
	noteViewport        viewport.Model
	spinner             spinner.Model
	loadingNote         bool
	textInput           textinput.Model
	passwordExists      bool
	creatingNewPassword bool
	newPasswordFocus    int
	pwConfirmTextInput  textinput.Model
	err                 error
}

func initialModel() model {
	items := []list.Item{
		item{title: "New note", desc: "Write a new encrypted note"},
	}

	textInput := textinput.New()
	textInput.Placeholder = "Password"
	textInput.EchoMode = textinput.EchoPassword
	textInput.Focus()

	s := spinner.New()
	s.Spinner = spinner.Points

	passwordExists, err := enotes.PasswordExists()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pwConfirmTextInput := textinput.New()
	pwConfirmTextInput.Placeholder = "Confirm Password"
	pwConfirmTextInput.EchoMode = textinput.EchoPassword

	m := model{
		chosen:             -1,
		list:               list.New(items, list.NewDefaultDelegate(), 0, 0),
		textInput:          textInput,
		noteViewport:       viewport.New(30, 20),
		spinner:            s,
		passwordExists:     passwordExists,
		pwConfirmTextInput: pwConfirmTextInput,
	}
	m.list.Title = "Notes"
	return m
}

func (m model) inNewPassword() bool {
	return !m.passwordExists
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
	return len(m.newNotePath) > 0
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

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, getDirFiles, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting || m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, nil
		}
	case editorFinishedMsg:
		m.editorActive = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.loadingNote = true
		if m.inNewNoteEditor() {
			m.newNotePath = ""
			m.resetChosen()
			return m, getDirFiles
		}
		item := m.list.SelectedItem().(fileItem)
		return m, openNote(item.file.Name(), m.password)
	case dirFilesMsg:
		itemsLen := len(m.list.Items())
		cmds := make([]tea.Cmd, len(msg.files))
		for i, file := range msg.files {
			var cmd tea.Cmd
			if itemsLen > i+1 {
				cmd = m.list.SetItem(i+1, fileItem{file})
			} else {
				cmd = m.list.InsertItem(i+1, fileItem{file})
			}
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case newPasswordMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.password = ""
		m.creatingNewPassword = false
		m.passwordExists = true
		m.textInput.SetValue("")
		return m, textinput.Blink
	case verifyPasswordMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if !msg.passwordsMatch {
			m.err = enotes.IncorrectPasswordError
			return m, nil
		}
		m.passwordVerified = true
		return m, nil
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, 100)
		m.noteViewport.Width = m.width
		m.noteViewport.Height = msg.Height

		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.inNewPassword() {
		return newPasswordUpdate(msg, m)
	}
	if m.inPassword() {
		return passwordUpdate(msg, m)
	}
	if !m.passwordVerified {
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
	if m.quitting {
		return ""
	}

	if m.err != nil {
		return m.err.Error() + "\n"
	}

	if m.inNewPassword() {
		return newPasswordView(m)
	}

	if m.inPassword() {
		return passwordView(m)
	}

	if !m.passwordVerified {
		return fmt.Sprintf("%s Verifying password\n", m.spinner.View())
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithMouseCellMotion())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
