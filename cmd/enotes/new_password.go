package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			return m, tea.Quit
		case "enter":
			switch m.newPasswordFocus {
			case 0:
				m.password = m.textInput.Value()
				m.newPasswordFocus += 1
				m.textInput.Blur()
				cmd := m.pwConfirmTextInput.Focus()
				return m, cmd
			case 1:
				confirmPw := m.pwConfirmTextInput.Value()
				m.pwConfirmTextInput.Blur()
				cmd := m.textInput.Focus()
				if m.password == confirmPw {
					return m, newPassword(m.password)
				} else {
					m.newPasswordsDontMatch = true
					return m, cmd
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.newPasswordFocus {
	case 0:
		m.textInput, cmd = m.textInput.Update(msg)
	case 1:
		var ti textinput.Model
		ti, cmd = m.pwConfirmTextInput.Update(msg)
		m.pwConfirmTextInput = &ti
	}
	return m, cmd
}

func newPasswordView(m model) string {
	if m.newPasswordsDontMatch {
		return "Password and confirm password didn't match"
	}
	return fmt.Sprintf(
		"New Password: %s\n\nConfirm password: %s\n\n%s",
		m.textInput.View(),
		m.pwConfirmTextInput.View(),
		"(esc to quit)",
	) + "\n"
}
