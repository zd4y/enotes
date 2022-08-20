package notes

import (
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
