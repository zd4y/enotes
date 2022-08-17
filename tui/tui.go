package tui

import (
	"fmt"
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

type model struct {
	list list.Model
	chosen int
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

func (m model) Init() tea.Cmd {
	return textinput.Blink
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
		m.resetChosen()
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
	return m, nil
}

func newNoteEditorUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	m.newNoteName = ""
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
	item := listItem.(item);
	return fmt.Sprintf("in note %s", item.title)
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

func Run() {
	items := []list.Item{
		item{title: "New note", desc: "Write a new encrypted note"},
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := model{chosen: -1, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
