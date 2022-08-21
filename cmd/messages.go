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
