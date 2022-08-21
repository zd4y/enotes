package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/zd4y/notes/pkg/notes"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct {
	path string
	err  error
	msg  tea.Msg
}

func openEditor(path string, callback func(error) tea.Msg) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		msg := callback(err)
		return editorFinishedMsg{path, err, msg}
	})
}

type noteEditedMsg struct {
	done func() error
	err  error
}

func editNote(notePath string, password string) tea.Cmd {
	return func() tea.Msg {
		tempNotePath, done, err := notes.EditNote(notePath, password)
		if err != nil {
			return noteEditedMsg{err: err}
		}
		return openEditor(tempNotePath, func(err error) tea.Msg {
			return noteEditedMsg{done, err}
		})()
	}
}

type verifyPasswordMsg struct {
	passwordsMatch bool
	err            error
}

func verifyPassword(password string) tea.Cmd {
	return func() tea.Msg {
		err := notes.VerifyPassword(password)
		return verifyPasswordMsg{err == nil, err}
	}
}

type openNoteMsg struct {
	note string
	err  error
}

func openNote(notePath string, password string) tea.Cmd {
	return func() tea.Msg {
		note, err := notes.OpenNote(notePath, password)
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
