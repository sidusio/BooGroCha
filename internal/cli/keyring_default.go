// +build !linux

package cli

import "fmt"

func hasKeyRingSupport() bool {
	return false
}

func getKeyRingPassword(key string) (string, error) {
	return "", fmt.Errorf("not implemented for current architecture")
}

func saveKeyRingPassword(key, pass string) error {
	return fmt.Errorf("not implemented for current architecture")
}
