/*
# File handling functions

F in OpenToReadF stands for... Do succesfully, or log to Fatal level and exit
*/
package file

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/bini"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type WebFile struct {
	url string
}

type File struct {
	filepath utils_types.FilePath
	file     *os.File
	lines    []string

	webfile *WebFile
}

func NewFile(filepath utils_types.FilePath) *File {
	return &File{filepath: filepath}
}

func NewWebFile(url string) *File {
	return &File{webfile: &WebFile{
		url: url,
	}}
}

func (f *File) GetFilepath() utils_types.FilePath { return f.filepath }

func (f *File) openToReadF() *File {
	logus.Log.Debug("opening file", utils_logus.FilePath(f.GetFilepath()))
	file, err := os.Open(string(f.filepath))
	f.file = file

	logus.Log.CheckPanic(err, "failed to open ", utils_logus.FilePath(f.filepath))
	return f
}

func (f *File) close() {
	f.file.Close()
}

func (f *File) ReadLines() ([]string, error) {

	if f.webfile != nil {
		res, err := http.Get(f.webfile.url)
		if err != nil {
			logus.Log.Error("error making http request: %s\n", typelog.OptError(err))
			return []string{}, err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			logus.Log.Error("client: could not read response body: %s\n", typelog.OptError(err))
			return []string{}, err
		}
		// fmt.Printf("client: response body: %s\n", resBody)

		str := string(resBody)
		return strings.Split(str, "\n"), nil
	}

	if bini.IsBini(f.filepath) {
		f.lines = bini.Dump(f.filepath)
		return f.lines, nil
	}

	f.openToReadF()
	defer f.close()

	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		f.lines = append(f.lines, scanner.Text())
	}
	return f.lines, nil
}

func (f *File) ScheduleToWrite(value ...string) {
	f.lines = append(f.lines, value...)
}

func (f *File) GetLines() []string {
	return f.lines
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
	logus.Log.CheckPanic(err, "failed to open ", utils_logus.FilePath(f.filepath))

	return f
}
func (f *File) writelnF(msg string) {
	_, err := f.file.WriteString(fmt.Sprintf("%v\n", msg))

	logus.Log.CheckPanic(err, "failed to write string to file")
}
