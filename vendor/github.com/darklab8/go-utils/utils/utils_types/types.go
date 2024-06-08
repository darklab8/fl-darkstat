package utils_types

import "path/filepath"

type FilePath string

func (f FilePath) ToString() string { return string(f) }

func (f FilePath) Base() FilePath { return FilePath(filepath.Base(string(f))) }

type RegExp string

type TemplateExpression string
