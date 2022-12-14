package tui

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/zd4y/enotes/enotes"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct {
	path string
	err  error
}

func openEditor(path string, callback func(error) error) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	editorCmd := strings.Split(editor, " ")
	editor, args := editorCmd[0], editorCmd[1:]
	args = append(args, path)
	c := exec.Command(editor, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		err = callback(err)
		return editorFinishedMsg{path, err}
	})
}

func createNote(noteName string, password string) tea.Cmd {
	return func() tea.Msg {
		tempNotePath, done, err := enotes.CreateNote(noteName, password)
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		return openEditor(tempNotePath, func(err error) error {
			if err != nil {
				return err
			}
			return done()
		})()
	}
}

func editNote(notePath string, password string) tea.Cmd {
	return func() tea.Msg {
		tempNotePath, done, err := enotes.EditNote(notePath, password)
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		return openEditor(tempNotePath, func(err error) error {
			if err != nil {
				return err
			}
			return done()
		})()
	}
}

type newPasswordMsg struct {
	err error
}

func newPassword(password string) tea.Cmd {
	return func() tea.Msg {
		err := enotes.NewPassword(password)
		return newPasswordMsg{err}
	}
}

type verifyPasswordMsg struct {
	passwordsMatch bool
	err            error
}

func verifyPassword(password string) tea.Cmd {
	return func() tea.Msg {
		err := enotes.VerifyPassword(password)
		return verifyPasswordMsg{err == nil, err}
	}
}

type openNoteMsg struct {
	note string
	err  error
}

func openNote(notePath string, password string) tea.Cmd {
	return func() tea.Msg {
		note, err := enotes.OpenNote(notePath, password)
		return openNoteMsg{note, err}
	}
}

type dirFilesMsg struct{ files []fs.FileInfo }

func getDirFiles() tea.Msg {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	return dirFilesMsg{files}
}
