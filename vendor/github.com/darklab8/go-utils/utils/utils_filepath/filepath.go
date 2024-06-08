package utils_filepath

import (
	"path/filepath"
	"strings"

	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func Join(paths ...utils_types.FilePath) utils_types.FilePath {
	return utils_types.FilePath(filepath.Join(utils.CompL(paths, func(path utils_types.FilePath) string { return string(path) })...))
}

func Dir(path utils_types.FilePath) utils_types.FilePath {
	return utils_types.FilePath(filepath.Dir(string(path)))
}

func Base(path utils_types.FilePath) utils_types.FilePath {
	return utils_types.FilePath(filepath.Base(string(path)))
}

func Contains(path utils_types.FilePath, subpath utils_types.FilePath) bool {
	return strings.Contains(path.ToString(), subpath.ToString())
}
