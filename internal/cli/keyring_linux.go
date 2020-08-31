package cli

import (
	"fmt"

	"github.com/jsipprell/keyctl"
)

func hasKeyRingSupport() bool {
	_, err := keyctl.SessionKeyring()
	return err == nil
}

func getKeyRingPassword(key string) (string, error) {

	keyring, err := keyctl.SessionKeyring()
	if err != nil {
		return "", err
	}
	pass, err := keyring.Search(key)
	notSaved := err != nil && err.Error() == "required key not available"
	if notSaved {
		return "", fmt.Errorf("not saved")
	} else if err != nil {
		return "", err
	}
	bytes, err := pass.Get()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func saveKeyRingPassword(key, pass string) error {
	// Setting a key to "" will result in crash, clear instead
	if pass == "" {
		return clearKeyringPassword(key)
	}
	keyring, err := keyctl.SessionKeyring()
	if err != nil {
		return err
	}
	_, err = keyring.Add(key, []byte(pass))
	return err
}

func clearKeyringPassword(key string) error {
	keyring, err := keyctl.SessionKeyring()
	if err != nil {
		return err
	}
	pass, err := keyring.Search(key)
	if err != nil {
		if err.Error() == "required key not available" {
			return nil
		}
		return err
	}
	return pass.Unlink()
}
