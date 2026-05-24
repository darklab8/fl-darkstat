package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_http"

	"github.com/darklab8/fl-darkstat/helpers/patch_disco"
)

const baseURL = "https://discoverygc.com/gameconfigpublic/"

var Log *typelog.Logger = typelog.NewLogger("discoapi", typelog.WithLogLevel(typelog.LEVEL_INFO))

func scrapeFileNames(logger *typelog.Logger) (error, []string) {
	res, err := utils_http.Get(baseURL)
	if res.StatusCode >= 400 {
		return errors.New(fmt.Sprintln("non positive status code in trying to parse file names, status=", res.StatusCode, " err=", err)), []string{}
	}
	if err != nil {
		return err, nil
	}
	resBody, err := io.ReadAll(res.Body)
	logger.Info(string(resBody))

	// Save file
	out, err := os.Create(filepath.Join(configs_path, "index.html"))
	if err != nil {
		return err, nil
	}
	defer out.Close()
	out.Write(resBody)
	// Save file end

	doc := soup.HTMLParse(string(resBody))
	links := doc.FindAll("a")

	var files []string
	for _, link := range links {
		href := link.Attrs()["href"]
		// Skip parent directory links and empty hrefs
		if href == "" || href == "../" || strings.HasPrefix(href, "?") || strings.HasPrefix(href, "/") {
			continue
		}
		files = append(files, href)
	}

	return nil, files
}

type ErrDownloadFile struct {
	status_code  int
	original_err error
}

func (e ErrDownloadFile) Error() string {
	return fmt.Sprintln(
		"error http request with non positive status code, status_code>=400, status_code=",
		e.status_code,
		" err=",
		e.original_err,
	)
}

func downloadFile(destDir, fileName string, urlInput string, inmemory bool) (err error, resBody []byte) {

	destPath := filepath.Join(destDir, fileName)

	url := baseURL + fileName
	if urlInput != "" {
		url = urlInput
	}

	resp, err := utils_http.Get(url)
	if resp.StatusCode >= 400 {
		return ErrDownloadFile{
			status_code:  resp.StatusCode,
			original_err: err,
		}, resBody
	}
	if err != nil {
		return err, resBody
	}

	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err, resBody
	}
	defer out.Close()

	if inmemory {
		resBody, err = io.ReadAll(resp.Body)
		out.Write(resBody)

	} else {
		_, err = io.Copy(out, resp.Body)
	}
	return err, resBody
}

var configs_path = "/data/gameconfigpublic"

func main() {
	os.Mkdir(configs_path, 0777)
	os.Mkdir("/data/forums", 0777)

	i := 0
	go func() {
		for {
			i++
			var run_errors []error
			downloaded_files := 0

			logger := Log.WithFields(typelog.Int("loop", i))
			logger.Infoln("Scraping file list from", baseURL, " version5")
			err, files := scrapeFileNames(logger)

			if logger.CheckError(err, "Error scraping file names") {
				time.Sleep(time.Minute * 3)
				continue
			}
			logger.Infoln("files=", files)

			logger.Info(fmt.Sprintf("Found %d files, downloading...\n", len(files)))
			for _, fileName := range files {
				logger := logger.WithFields(typelog.Any("filename", fileName))
				logger.Info("Downloading: %s\n")
				err, _ := downloadFile(configs_path, fileName, "", false)
				if err != nil {
					var myErr ErrDownloadFile
					if errors.As(err, &myErr) {
						if myErr.status_code == 404 {
							logger.CheckWarn(err, "downloading file failed with not found error")
						} else {
							logger.CheckError(err, "downloading file failed")
							run_errors = append(run_errors, err)
						}
					} else {
						logger.CheckError(err, "downloading file failed")
					}
				} else {
					downloaded_files++
				}
			}

			err, data := downloadFile("/data", "patchlist.xml", "https://patch.discoverygc.com/patchlist.xml", false)
			if logger.CheckError(err, "https://patch.discoverygc.com/patchlist.xml Error downloading patchlist.xml") {
				run_errors = append(run_errors, err)
			}

			err, data = downloadFile("/data", "forums/base_admin.php", "https://discoverygc.com/forums/base_admin.php?action=getjson", true)
			if logger.CheckError(err, "base_admin.php5 Error downloading forums/base_admin.php") {
				run_errors = append(run_errors, err)
			}

			if len(data) < 1000 {
				err := errors.New(fmt.Sprintln("base_admin.php5 is too small (showing content). time=", time.Now(), " len=", len(data), string(data)))
				if logger.CheckError(err, "base_adminphp5 error") {
					run_errors = append(run_errors, err)
				}
			} else {
				logger.Infoln("base_admin downloaded succesfully. time=", time.Now(), " len=", len(data))
			}

			unmarshaled := make(map[string]any)
			err = json.Unmarshal(data, &unmarshaled)
			if logger.CheckErrorln(err, "base_admin.php5 failed to unmarshal its json (showing no content) ", time.Now(), " len=", len(data)) {
				err = os.WriteFile("/data/errored_base_admin.json", data, os.FileMode(0644))
				if Log.CheckError(err, "base_admin.php5 failed to write errored data to file") {
					run_errors = append(run_errors, err)
				}

			}

			err, patchlistdata := downloadFile("/data", "patchlist.xml", "https://patch.discoverygc.com/patchlist.xml", true)
			if logger.CheckError(err, "Error downloading patchlist.xml") {
				run_errors = append(run_errors, err)
			} else {
				Log.Info("downloaded patchlist.xml succesfully")
			}

			patches := patch_disco.ParseForPatches(patch_disco.DiscoveryUrl, patchlistdata)
			if len(patches) == 0 {
				logger.Error("Parsed zero patches from patchlist.xml")
			} else {
				for _, patch := range patches[len(patches)-10:] {

					if FileExists(filepath.Join("/data", patch.Filename)) {
						Log.Warnln("patch already exists, skipping downloading it, filename=", patch.Filename)
						continue
					} else {
						Log.Infoln("patch does not exist, downloading, filename=", patch.Filename)
					}
					err, _ := downloadFile("/data", patch.Filename, patch.Url, false)
					if logger.CheckErrorln(err, "Error downloading patch=", patch.Filename) {
						run_errors = append(run_errors, err)
					}

				}
			}

			logger.Info("All downloads complete5.")
			logger = logger.WithFields(
				typelog.Int("files_count", len(files)),
				typelog.Int("run_errors_count", len(run_errors)),
			)
			if len(run_errors) == 0 && len(files) > 0 {
				logger.Info("run was succesful")
			} else {
				logger.Info("run was failure")
			}

			time.Sleep(time.Minute * 3)
		}
	}()

	port := flag.String("p", "8000", "port to serve on")
	directory := "/data"
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(directory)))

	Log.Info(fmt.Sprintf("Serving %s on HTTP port: %s\n", directory, *port))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
