package file

import (
	"io"
	"os"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_http"
)

func (f *File) ReadBytes() ([]byte, error) {
	if f.webfile != nil {
		res, err := utils_http.Get(f.webfile.url)
		if logus.Log.CheckError(err, "error making http request: %s\n", typelog.OptError(err)) {
			return []byte{}, err
		}

		resBody, err := io.ReadAll(res.Body)
		if logus.Log.CheckError(err, "client: could not read response body: %s\n", typelog.OptError(err)) {
			return []byte{}, err
		}

		return resBody, nil
	}

	resBody, err := os.ReadFile(f.filepath.ToString())
	logus.Log.CheckError(err, "client: could not read os.ReadFile body: %s\n", typelog.OptError(err))
	return resBody, err
}
