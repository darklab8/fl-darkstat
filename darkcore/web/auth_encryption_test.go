package web

import (
	"testing"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/stretchr/testify/assert"
)

func TestEncr(t *testing.T) {
	key := "passphrasewhichneedstobe32bytes!"
	phrase := "myphrase"

	encrypted, err := encrypt([]byte(phrase), key)
	logus.Log.CheckPanic(err, "failed encrypting")

	decrypted, err := decrypt(encrypted, key)
	logus.Log.CheckPanic(err, "failed decryption")

	assert.Equal(t, phrase, string(decrypted))

}
