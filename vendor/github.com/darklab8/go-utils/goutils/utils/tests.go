package utils

import (
	"os"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
)

func RegenerativeTest(callback func() error) error {
	if os.Getenv("DARK_TEST_REGENERATE") != "true" {
		utils_logus.Log.Debug("Skipping test data regenerative code")
		return nil
	}

	return callback()
}
