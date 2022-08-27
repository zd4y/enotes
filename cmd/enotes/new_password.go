package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		msg := msg.String()
		switch msg {
		case "esc":
			m.quitting = true
			return m, nil
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
					m.creatingNewPassword = true
					return m, tea.Batch(newPassword(m.password), cmd)
				} else {
					m.err = errors.New("Password and confirm password didn't match")
					return m, nil
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
		m.pwConfirmTextInput = ti
	}
	return m, cmd
}

func newPasswordView(m model) string {
	if m.creatingNewPassword {
		return fmt.Sprintf("%s Creating password\n", m.spinner.View())
	}

	return fmt.Sprintf(
		"New Password: %s\n\nConfirm password: %s\n\n%s\n",
		m.textInput.View(),
		m.pwConfirmTextInput.View(),
		"(esc to quit)",
	)
}
