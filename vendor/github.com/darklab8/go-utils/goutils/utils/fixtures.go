package utils

import "os"

func FixtureDevEnv() bool {
	return os.Getenv("DEV_ENV") == "true"
}
