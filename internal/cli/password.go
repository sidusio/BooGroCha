package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

const KeyringName = ApplicationName + "-password"

func getPassword() string {
	var password string

	if !hasKeyRingSupport() {
		if viper.GetString("chalmers.pass") == "" {
			fmt.Println("No password set in config. You can set it with 'bgc config set pass'")
		}

		// Prompt for password if it's not set in config or the cid flag is specified
		if viper.GetString("chalmers.pass") == "" || BgcCmd.Flag("cid").Value.String() != "" {
			password = promptForPassword()
		} else {
			bytes, err := base64.StdEncoding.DecodeString(viper.GetString("chalmers.pass"))
			if err != nil {
				fmt.Printf("Failed to read password: %s\n", err.Error())
				fmt.Printf("Try to reset or unset your password with 'gbc config set pass'\n")
				os.Exit(1)
			}
			password = string(bytes)
		}
	} else {
		key, err := getKeyRingPassword(KeyringName)
		notSaved := err != nil && err.Error() == "not saved"
		if notSaved {
			fmt.Println("No password set. You can set it securely with 'bgc config set pass'")
		}
		if err != nil && !notSaved {
			fmt.Printf("Failed to read password: %s\n", err.Error())
			fmt.Printf("Try to reset or unset your password with 'gbc config set pass'\n")
			os.Exit(1)
		}
		if err != nil || BgcCmd.Flag("cid").Value.String() != "" {
			password = promptForPassword()
		} else {
			password = key
		}
	}
	return password
}

func promptForPassword() string {
	fmt.Print("Enter Password: ")
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("failed to read password")
		os.Exit(1)
	}
	return strings.TrimSpace(string(bytes))
}

func getSavePassword() func(string) error {
	if !hasKeyRingSupport() {
		return nil
	}
	return func(password string) error {
		return saveKeyRingPassword(KeyringName, password)
	}
}
