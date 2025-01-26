package web

import (
	"encoding/json"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
)

type AuthPermission struct {
	ExpirationDate time.Time
	IsValid        bool
}

func NewTempusToken() string {
	data := AuthPermission{ExpirationDate: time.Now().Add(12 * time.Hour), IsValid: true}
	bytes_data, err := json.Marshal(data)
	logus.Log.CheckError(err, "failed to marshal auth permission")

	encrypted, err := encrypt(bytes_data, settings.Env.Secret)
	logus.Log.CheckError(err, "failed encrypting")
	return encrypted
}

func IsTempusValid(encrypted string) bool {
	decrypted, err := decrypt(encrypted, settings.Env.Secret)
	if logus.Log.CheckError(err, "failed decryption") {
		return false
	}

	var data2 AuthPermission
	err = json.Unmarshal(decrypted, &data2)
	if logus.Log.CheckError(err, "failed to unmarshal") {
		return false
	}

	if !data2.IsValid {
		return false
	}
	if !time.Now().Before(data2.ExpirationDate) {
		return false
	}

	return true
}
