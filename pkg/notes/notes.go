package notes

import (
	"bytes"
	"io"
	"os"

	"filippo.io/age"
)

const passwordFileName = ".notes-password.age"

func VerifyPassword(password string) error {
	passwordFile, err := os.Open(passwordFileName)
	if err != nil {
		return err
	}
	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		return err
	}
	_, err = age.Decrypt(passwordFile, identity)
	return err
}

func OpenNote(path string, password string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		return "", err
	}

	r, err := age.Decrypt(file, identity)
	if err != nil {
		return "", err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return "", err
	}

	return out.String(), nil
}
