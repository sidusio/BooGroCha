package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	c := Credentials{
		CID:      "Emil",
		Password: "Emil's pass",
	}

	secret := []byte("00000secret00000")

	cipher, err := Encrypt(c, secret)
	assert.NoError(t, err)

	cx, err := Decrypt(cipher, secret)
	assert.NoError(t, err)

	assert.Equal(t, c, cx)
}
