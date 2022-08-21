package enotes

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"filippo.io/age"
)

const passwordFileName = ".enotes-password.age"

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
	bytes, err := decrypt(path, password)
	if err != nil {
		return "", err
	}
	return bytes.String(), nil
}

func CreateNote(path string, password string) (string, func() error, error) {
	tempFileName, done, err := useTempFile(
		path,
		func(tempFile *os.File) error { return nil },
	)
	if err != nil {
		return "", nil, err
	}

	return tempFileName, func() error {
		err := encryptFile(tempFileName, path, password)
		if err != nil {
			return err
		}
		return done()
	}, nil
}

func EditNote(path string, password string) (string, func() error, error) {
	noteBytes, err := decrypt(path, password)
	if err != nil {
		return "", nil, err
	}

	tempFileName, done, err := useTempFile(
		path,
		func(tempFile *os.File) error {
			_, err = tempFile.Write(noteBytes.Bytes())
			return err
		},
	)
	if err != nil {
		return "", nil, err
	}

	return tempFileName, func() error {
		err := encryptFile(tempFileName, path, password)
		if err != nil {
			return err
		}
		return done()
	}, err
}

func decrypt(path string, password string) (*bytes.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		return nil, err
	}

	r, err := age.Decrypt(file, identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}
	return out, nil
}

func encryptFile(srcPath string, dstPath string, password string) error {
	srcBytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	recipient, err := age.NewScryptRecipient(password)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	w, err := age.Encrypt(dstFile, recipient)
	if err != nil {
		return err
	}

	_, err = w.Write(srcBytes)
	if err != nil {
		return err
	}

	return w.Close()
}

func useTempFile(prefix string, manipulateTempFile func(*os.File) error) (string, func() error, error) {
	tempFile, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", nil, err
	}

	err = manipulateTempFile(tempFile)
	if err != nil {
		return "", nil, err
	}

	err = tempFile.Close()
	if err != nil {
		return "", nil, err
	}

	tempFileName := tempFile.Name()

	return tempFileName, func() error {
		return os.Remove(tempFileName)
	}, nil
}
