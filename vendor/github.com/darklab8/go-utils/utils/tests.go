package utils

import (
	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

func RegenerativeTest(callback func() error) error {
	if !utils_settings.Envs.AreTestsRegenerating {
		utils_logus.Log.Debug("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
