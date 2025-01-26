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

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/bini"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type WebFile struct {
	url string
}

type File struct {
	filepath utils_types.FilePath
	file     *os.File
	lines    []string

	// failback files escape fate of being written in
	IsFailback bool

	webfile *WebFile
}

func NewMemoryFile(lines []string) *File {
	return &File{lines: lines}
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

	if len(f.lines) > 0 {
		lines := f.lines
		f.lines = []string{}
		return lines, nil
	}

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
		bini_lines := bini.Dump(f.filepath)
		return bini_lines, nil
	}

	f.openToReadF()
	defer f.close()

	scanner := bufio.NewScanner(f.file)

	bufio_lines := []string{}
	for scanner.Scan() {
		bufio_lines = append(bufio_lines, scanner.Text())
	}
	return bufio_lines, nil
}

func (f *File) ScheduleToWrite(value ...string) {
	f.lines = append(f.lines, value...)
}

func (f *File) WriteLines() {
	if f.IsFailback {
		// This feature is not working in full capacity for some reason ;) not getting skipped for some reason
		logus.Log.Warn("file is taken from fallback, writing is skipped",
			typelog.Any("filename", f.filepath.ToString()),
		)
		return
	}

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
