package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
)

func Encrypt(credentials Credentials, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	plaintext, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) (Credentials, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return Credentials{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return Credentials{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return Credentials{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	credentials := Credentials{}
	err = json.Unmarshal(plaintext, &credentials)

	return credentials, err
}
