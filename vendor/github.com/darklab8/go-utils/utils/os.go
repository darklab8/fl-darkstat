package utils

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func GetCurrentFile() utils_types.FilePath {
	_, filename, _, _ := runtime.Caller(1)
	return utils_types.FilePath(filename)
}

func GetCurrentFolder() utils_types.FilePath {
	_, filename, _, _ := runtime.Caller(1)
	directory := filepath.Dir(filename)
	return utils_types.FilePath(directory)
}

func GetProjectDir() utils_types.FilePath {
	path, err := os.Getwd()
	if folder_override, ok := os.LookupEnv("AUTOGIT_PROJECT_FOLDER"); ok {
		path = folder_override
	}
	utils_logus.Log.CheckPanic(err, "unable to get workdir")
	return utils_types.FilePath(path)
}

func GetCurrrentChildFolder(folder_name string) utils_types.FilePath {
	_, filename, _, _ := runtime.Caller(2)
	directory := filepath.Dir(filename)
	test_directory := filepath.Join(directory, folder_name)
	return utils_types.FilePath(test_directory)
}

func GetCurrrentTestFolder() utils_types.FilePath {
	return GetCurrrentChildFolder("testdata")
}

func GetCurrrentTempFolder() utils_types.FilePath {
	return GetCurrrentChildFolder("tempdata")
}
