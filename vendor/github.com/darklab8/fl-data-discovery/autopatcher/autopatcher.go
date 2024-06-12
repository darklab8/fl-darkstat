package autopatcher

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type RequestResp struct {
	Body       []byte
	StatusCode int
}

func Request(url string) RequestResp {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	// fmt.Printf("client: response body: %s\n", resBody)
	return RequestResp{
		Body:       resBody,
		StatusCode: res.StatusCode,
	}
}

func downloadFile(filepath string, url string) (err error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(fpath string) bool {
	if _, err := os.Stat(fpath); err == nil {
		// path/to/whatever exists
		return true
	}
	return false
}

/*
Unzip is copy paste from
https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file
https://stackoverflow.com/a/24792688
*/
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0777)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0777)
		} else {
			os.MkdirAll(filepath.Dir(path), 0777)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

type PatcherData struct {
	XMLName   xml.Name `xml:"PatcherData"`
	Text      string   `xml:",chardata"`
	PatchList struct {
		Text  string `xml:",chardata"`
		Patch []struct {
			Text    string `xml:",chardata"`
			URL     string `xml:"url,attr"`
			Md5hash string `xml:"md5hash,attr"`
		} `xml:"patch"`
	} `xml:"PatchList"`
}

type Patch struct {
	Filename string
	Url      string
	Hash     PatchHash
	Name     string
}

func (patch Patch) GetFilepath() string {
	return filepath.Join("patches", patch.Filename)
}

func (patch Patch) GetFolderPath() string {
	return filepath.Join("patches", strings.ReplaceAll(patch.Filename, ".zip", ""))
}

func parseForPatches(discovery_url string, body []byte) []Patch {
	var patches []Patch
	var Page PatcherData
	xml.Unmarshal(body, &Page)

	for _, patch := range Page.PatchList.Patch {
		// fmt.Println(patch)
		patches = append(patches, Patch{
			Filename: patch.URL,
			Url:      discovery_url + patch.URL,
			Hash:     PatchHash(patch.Md5hash),
			Name:     patch.Text,
		})
	}
	return patches
}

func downloadPatch(patch Patch) {
	os.MkdirAll("patches", 0777)
	if fileExists(patch.GetFilepath()) {
		fmt.Println("fpath already eixsts, fpath=", patch.GetFilepath())
		return
	}

	err := downloadFile(patch.GetFilepath(), patch.Url)
	if err != nil {
		fmt.Println("not able to download url", err, patch, "url=", patch.Url)
		os.Exit(1)
	}
	fmt.Println("downloaded file", patch.GetFilepath())
}

type File struct {
	filepath_ string
}

func (f File) GetPath() string {
	return f.filepath_
}

func (f File) GetRelPathTo(root string) string {
	path := strings.ReplaceAll(f.filepath_, root+PATH_SEPARATOR, "")
	return path
}

func (f File) GetLowerPath() string {
	return strings.ToLower(f.filepath_)
}

func NewFile(path string) File {
	f := File{filepath_: path}
	return f
}

type Filesystem struct {
	Files         []File
	LowerMapFiles map[string]File

	Folders         []File
	LowerMapFolders map[string]File
}

func ScanCaseInsensitiveFS(fs_path string) Filesystem {
	myfs := Filesystem{
		LowerMapFiles:   make(map[string]File),
		LowerMapFolders: make(map[string]File),
	}
	filepath.WalkDir(fs_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			file := NewFile(path)
			myfs.Files = append(myfs.Files, file)
			myfs.LowerMapFiles[file.GetLowerPath()] = file
		} else {
			folder := NewFile(path)
			myfs.Folders = append(myfs.Folders, folder)
			myfs.LowerMapFolders[folder.GetLowerPath()] = folder
		}
		return nil

	})
	return myfs
}

func WriteToFile(path string, content []byte) {
	destination, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintln("os.Create:", err))
	}
	defer destination.Close()

	_, err = destination.Write(content)
	if err != nil {
		panic(fmt.Sprintln("failed to write file", err))
	}
}

type BadassRoot struct {
	XMLName      xml.Name `xml:"BadassRoot"`
	Text         string   `xml:",chardata"`
	Xsd          string   `xml:"xsd,attr"`
	Xsi          string   `xml:"xsi,attr"`
	PatchHistory struct {
		Text  string   `xml:",chardata"`
		Patch []string `xml:"Patch"`
	} `xml:"PatchHistory"`
}

type PatchHash string

func readLauncherConfig() map[PatchHash]string {
	body, err := os.ReadFile("launcherconfig.xml")
	if err != nil {
		panic(fmt.Sprintln("failed to read launcherconfig", err))
	}

	var patches map[PatchHash]string = make(map[PatchHash]string)
	var Page BadassRoot
	xml.Unmarshal(body, &Page)

	for _, patch := range Page.PatchHistory.Patch {
		patches[PatchHash(patch)] = ""
	}
	return patches
}

/*
AdjustFoldersInPatch replaces folder names in path to relevant case sensitive folders
And creates them if they don't exist
*/
func AdjustFoldersInPath(relative_patch_filepath string, freelancer_folder Filesystem) string {
	patch_chain_folders := strings.Split(relative_patch_filepath, PATH_SEPARATOR)
	patch_chain_folders = patch_chain_folders[:len(patch_chain_folders)-1] // minus filename

	for i := 1; i <= len(patch_chain_folders); i++ {
		folders_chain := patch_chain_folders[:i]
		patch_target_folder := filepath.Join(folders_chain...)

		if folder, found := freelancer_folder.LowerMapFolders[strings.ToLower(patch_target_folder)]; found {
			relative_patch_filepath = strings.ReplaceAll(relative_patch_filepath, patch_target_folder, folder.GetPath())

			// refreshing
			patch_chain_folders = strings.Split(relative_patch_filepath, PATH_SEPARATOR)
			patch_chain_folders = patch_chain_folders[:len(patch_chain_folders)-1] // minus filename
		} else {
			err := os.MkdirAll(patch_target_folder, os.ModePerm)
			if err != nil {
				panic(fmt.Sprintln("failed creating mkdirall", err))
			}
		}
	}

	return relative_patch_filepath
}

var PATH_SEPARATOR = ""

func init() {
	if runtime.GOOS == "windows" {
		PATH_SEPARATOR = "\\"
	} else {
		PATH_SEPARATOR = "/"
	}
}

func RunAutopatcher() {
	discovery_url := "https://patch.discoverygc.com/"
	discovery_path_url := discovery_url + "patchlist.xml"
	resp := Request(discovery_path_url)

	patches := parseForPatches(discovery_url, resp.Body)

	patchhistory := readLauncherConfig()

	for _, patch := range patches {

		if _, found := patchhistory[patch.Hash]; found {
			fmt.Println("patch is already installed", patch)
			continue
		}

		downloadPatch(patch)

		patch_body, _ := os.ReadFile(patch.GetFilepath())
		hash := md5.Sum(patch_body)
		md5_result := hex.EncodeToString(hash[:])
		fmt.Println("md5_result=", md5_result)

		if md5_result != strings.ToLower(string(patch.Hash)) {
			panic(fmt.Sprintln("md5 hash sum is not matching", "expected=", patch.Hash, " but bound=", md5_result))
		}

		// md5, err := checksum.MD5sum(patch.GetFilepath())
		// fmt.Println(md5)
		// fmt.Println("md5_result=", md5, err)

		Unzip(patch.GetFilepath(), patch.GetFolderPath())

		freelancer_folder := ScanCaseInsensitiveFS(".")
		patch_folder := ScanCaseInsensitiveFS(patch.GetFolderPath())

		for _, file := range patch_folder.Files {

			content, err := os.ReadFile(file.GetPath())
			if err != nil {
				panic(fmt.Sprintln("failed to read file", err))
			}

			relative_patch_filepath := file.GetRelPathTo(patch.GetFolderPath())

			if strings.Contains(relative_patch_filepath, ".gitignore") {
				continue
			}

			if freelancer_path, file_exists := freelancer_folder.LowerMapFiles[strings.ToLower(relative_patch_filepath)]; file_exists {
				os.Remove(freelancer_path.GetPath())
			}

			relative_patch_filepath = AdjustFoldersInPath(relative_patch_filepath, freelancer_folder)

			WriteToFile(relative_patch_filepath, content)

		}

		os.RemoveAll(patch.GetFolderPath())

		fmt.Println("applied patch", patch)
		patch_marshaled, _ := json.Marshal(patch)
		os.WriteFile(AutopatherFilename, patch_marshaled, 0666)
	}
}

const AutopatherFilename = "autopatcher.latest_patch.json"
