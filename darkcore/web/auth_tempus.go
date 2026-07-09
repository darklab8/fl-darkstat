package web

import (
	"encoding/json"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-utils/typelog"
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

func IsTempusValid(encrypted string, logger *typelog.Logger) bool {
	logger = logger.WithFields(typelog.String("encrypted_str", encrypted))
	decrypted, err := decrypt(encrypted, settings.Env.Secret)
	if logger.CheckWarn(err, "failed decryption") {
		return false
	}

	var data2 AuthPermission
	err = json.Unmarshal(decrypted, &data2)
	if logger.CheckError(err, "failed to unmarshal") {
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
