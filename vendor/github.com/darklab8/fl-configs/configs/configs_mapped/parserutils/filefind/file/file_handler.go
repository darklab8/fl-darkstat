/*
# File handling functions

F in OpenToReadF stands for... Do succesfully, or log to Fatal level and exit
*/
package file

import (
	"bufio"
	"fmt"
	"os"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/bini"
	"github.com/darklab8/fl-configs/configs/settings/logus"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type File struct {
	filepath utils_types.FilePath
	file     *os.File
	lines    []string
}

func NewFile(filepath utils_types.FilePath) *File {
	return &File{filepath: filepath}
}

func (f *File) GetFilepath() utils_types.FilePath { return f.filepath }

func (f *File) openToReadF() *File {
	logus.Log.Debug("opening file", utils_logus.FilePath(f.GetFilepath()))
	file, err := os.Open(string(f.filepath))
	f.file = file

	logus.Log.CheckFatal(err, "failed to open ", utils_logus.FilePath(f.filepath))
	return f
}

func (f *File) close() {
	f.file.Close()
}

func (f *File) ReadLines() []string {
	if bini.IsBini(f.filepath) {
		return bini.Dump(f.filepath)
	}

	f.openToReadF()
	defer f.close()

	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		f.lines = append(f.lines, scanner.Text())
	}
	return f.lines
}

func (f *File) ScheduleToWrite(value ...string) {
	f.lines = append(f.lines, value...)
}

func (f *File) WriteLines() {
	f.createToWriteF()
	defer f.close()

	for _, line := range f.lines {
		f.writelnF(line)
	}
}

func (f *File) createToWriteF() *File {
	file, err := os.Create(string(f.filepath))
	f.file = file
	logus.Log.CheckFatal(err, "failed to open ", utils_logus.FilePath(f.filepath))

	return f
}
func (f *File) writelnF(msg string) {
	_, err := f.file.WriteString(fmt.Sprintf("%v\n", msg))

	logus.Log.CheckFatal(err, "failed to write string to file")
}
