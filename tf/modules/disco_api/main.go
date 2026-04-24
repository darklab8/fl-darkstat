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
	"github.com/darklab8/go-utils/examples/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_http"
)

const baseURL = "https://discoverygc.com/gameconfigpublic/"

func scrapeFileNames() ([]string, error) {
	res, err := utils_http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		logus.Log.CheckError(err, "error http request with non positive status code, status_code>=400", typelog.Any("status_code", res.StatusCode))
		return []string{}, errors.New(fmt.Sprintln("non positive status code in trying to parse file names, status=", res.StatusCode))
	}
	resBody, err := io.ReadAll(res.Body)
	fmt.Println(string(resBody))

	// Save file
	out, err := os.Create(filepath.Join(configs_path, "index.html"))
	if err != nil {
		return nil, err
	}
	defer out.Close()
	out.Write(resBody)
	// Save file end

	soup.HTMLParse(string(resBody))

	resp, err := soup.Get(baseURL)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(resp)
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

	return files, nil
}

func downloadFile(destDir, fileName string, urlInput string, inmemory bool) (err error, resBody []byte) {

	destPath := filepath.Join(destDir, fileName)

	url := baseURL + fileName
	if urlInput != "" {
		url = urlInput
	}

	resp, err := utils_http.Get(url)
	if err != nil {
		return err, resBody
	}
	if resp.StatusCode >= 400 {
		logus.Log.CheckError(err, "error http request with non positive status code, status_code>=400", typelog.Any("status_code", resp.StatusCode))
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

	go func() {
		for {
			log.Println("Scraping file list from", baseURL, " version5")
			files, err := scrapeFileNames()
			if err != nil {
				log.Printf("Error scraping file names: %v\n", err)
				time.Sleep(time.Minute * 3)
				continue
			}
			fmt.Println("files=", files)

			log.Printf("Found %d files, downloading...\n", len(files))
			for _, fileName := range files {
				log.Printf("Downloading: %s\n", fileName)
				if err, _ := downloadFile(configs_path, fileName, "", false); err != nil {
					log.Printf("Error downloading %s: %v\n", fileName, err)
				}
			}

			err, data := downloadFile("/data", "forums/base_admin.php", "https://discoverygc.com/forums/base_admin.php?action=getjson", true)
			if err != nil {
				log.Printf("ERROR base_admin.php5 Error downloading %s: %v\n", "forums/base_admin.php", err)
			}

			if len(data) < 1000 {
				log.Println("ERROR base_admin.php5 is too small (showing content). time=", time.Now(), " len=", len(data), string(data))
				// log.Println("base_admin.php4 is too small. (showing with content) time=", time.Now(), " len=", len(data), " content=", string(data))
			} else {
				log.Println("base_admin downloaded succesfully. time=", time.Now(), " len=", len(data))
			}

			unmarshaled := make(map[string]any)
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				log.Println("ERROR base_admin.php5 failed to unmarshal its json (showing no content) ", time.Now(), " len=", len(data))
				// log.Println("base_admin.php4 failed to unmarshal its json (showing with content)", time.Now(), " len=", len(data), " content=", string(data))
				err = os.WriteFile("/data/errored_base_admin.json", data, os.FileMode(0644))
				if err != nil {
					log.Println("ERRPR base_admin.php5 failed to write errored data to file")
				}
			}

			log.Println("All downloads complete5.")
			time.Sleep(time.Minute * 3)
		}
	}()

	port := flag.String("p", "8000", "port to serve on")
	directory := "/data"
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(directory)))

	log.Printf("Serving %s on HTTP port: %s\n", directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
