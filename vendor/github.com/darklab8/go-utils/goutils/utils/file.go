package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type file struct {
	filepath utils_types.FilePath

	file  *os.File
	lines []string
}

type FileRead struct {
	file
}

func NewReadFile(filepath utils_types.FilePath, callback func(*FileRead)) {
	f := &FileRead{file{filepath: filepath}}

	file, err := os.Open(string(f.filepath))
	f.file.file = file

	utils_logus.Log.CheckFatal(err, "failed to open", utils_logus.FilePath(f.filepath))
	defer f.file.file.Close()

	callback(f)
}

func (f *FileRead) ReadLines() []string {

	scanner := bufio.NewScanner(f.file.file)
	f.lines = []string{}
	for scanner.Scan() {
		f.lines = append(f.lines, scanner.Text())
	}
	return f.lines
}

type FileWrite struct {
	file
}

func NewWriteFile(filepath utils_types.FilePath, callback func(*FileWrite)) {
	f := &FileWrite{file{filepath: filepath}}

	file, err := os.Create(string(f.filepath))
	f.file.file = file
	utils_logus.Log.CheckFatal(err, "failed to open ", utils_logus.FilePath(f.filepath))
	defer f.file.file.Close()
	callback(f)
}

func (f *FileWrite) WritelnF(msg string) {
	_, err := f.file.file.WriteString(fmt.Sprintf("%v\n", msg))

	utils_logus.Log.CheckFatal(err, "failed to write string to file")
}

func FileExists(filename utils_types.FilePath) bool {
	info, err := os.Stat(string(filename))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
