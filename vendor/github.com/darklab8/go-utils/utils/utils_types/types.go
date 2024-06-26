package utils_types

import "path/filepath"

type FilePath string

func (f FilePath) ToString() string { return string(f) }

func (f FilePath) Base() FilePath { return FilePath(filepath.Base(string(f))) }

func (f FilePath) Dir() FilePath { return FilePath(filepath.Dir(string(f))) }

func (f FilePath) Join(paths ...string) FilePath {
	paths = append([]string{string(f)}, paths...)
	return FilePath(filepath.Join(paths...))
}

type RegExp string

type TemplateExpression string
