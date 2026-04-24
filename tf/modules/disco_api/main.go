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

func downloadFile(destDir, fileName string, urlInput string) error {

	destPath := filepath.Join(destDir, fileName)

	url := baseURL + fileName
	if urlInput != "" {
		url = urlInput
	}

	resp, err := utils_http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		logus.Log.CheckError(err, "error http request with non positive status code, status_code>=400", typelog.Any("status_code", resp.StatusCode))
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

var configs_path = "/data/gameconfigpublic"

func main() {
	os.Mkdir(configs_path, 0777)
	os.Mkdir("/data/forums", 0777)

	go func() {
		for {
			log.Println("Scraping file list from", baseURL)
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
				if err := downloadFile(configs_path, fileName, ""); err != nil {
					log.Printf("Error downloading %s: %v\n", fileName, err)
				}
			}

			if err := downloadFile("/data", "forums/base_admin.php", "https://discoverygc.com/forums/base_admin.php?action=getjson"); err != nil {
				log.Printf("Error downloading %s: %v\n", "forums/base_admin.php", err)
			}

			data, err := os.ReadFile("/data/forums/base_admin.php")
			if err != nil {
				log.Println("Error reading base_admin.php file to validate it", time.Now())
			} else {
				if len(data) < 1000 {
					log.Println("base_admin.php is too small (showing no content). time=", time.Now(), " len=", len(data))
					log.Println("base_admin.php is too small. (showing with content) time=", time.Now(), " len=", len(data), " content=", string(data))
				} else {
					log.Println("base_admin.php is has size. time=", time.Now(), " len=", len(data))
				}
			}

			unmarshaled := make(map[any]any)
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				log.Println("base_admin.php failed to unmarshal its json (showing no content) ", time.Now(), " len=", len(data))
				log.Println("base_admin.php failed to unmarshal its json (showing with content)", time.Now(), " len=", len(data), " content=", string(data))
			}
			log.Println("All downloads complete.")
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
