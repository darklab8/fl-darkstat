package patch_disco

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_http"
)

type RequestResp struct {
	Body       []byte
	StatusCode int
}

var Log *typelog.Logger = typelog.NewLogger(
	"disco_patch",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)

func Request(url string) (RequestResp, error) {
	res, err := utils_http.Get(url)

	if res.StatusCode >= 400 {
		return RequestResp{}, errors.New(fmt.Sprintln("status code is bad, status=", res.StatusCode, err))
	}

	if Log.CheckError(err, "client: error making http request to url", typelog.Any("url", url)) {
		return RequestResp{}, err
	}

	Log.Info("client: got response!", typelog.Any("status_code", res.StatusCode), typelog.Any("url", url))

	resBody, err := io.ReadAll(res.Body)

	if Log.CheckError(err, "client: could not read response body", typelog.Any("url", url)) {
		return RequestResp{}, err
	}
	return RequestResp{
		Body:       resBody,
		StatusCode: res.StatusCode,
	}, nil
}

func DownloadFile(filepath string, url string) (err error) {
	resp, err := utils_http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		Log.CheckError(err, "failed to close body")
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading file bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		Log.CheckError(err, "failed to close file to write")
	}()

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

	_ = os.MkdirAll(dest, 0777)

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
			_ = os.MkdirAll(path, 0777)
		} else {
			_ = os.MkdirAll(filepath.Dir(path), 0777)
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

func ParseForPatches(discovery_url string, body []byte) []Patch {
	var patches []Patch
	var Page PatcherData
	err := xml.Unmarshal(body, &Page)
	Log.CheckError(err, "failed unmarshaling xml for parse for patches")

	for _, patch := range Page.PatchList.Patch {
		patches = append(patches, Patch{
			Filename: patch.URL,
			Url:      discovery_url + patch.URL,
			Hash:     PatchHash(patch.Md5hash),
			Name:     patch.Text,
		})
	}
	return patches
}

func downloadPatch(patch Patch) error {
	_ = os.MkdirAll("patches", 0777)
	if fileExists(patch.GetFilepath()) {
		return errors.New(fmt.Sprintln("fpath already eixsts, fpath=", patch.GetFilepath()))
	}

	err := DownloadFile(patch.GetFilepath(), patch.Url)
	if err != nil {
		err_msg := fmt.Sprintln("not able to download url", err, patch, "url=", patch.Url)
		return errors.New(err_msg)
	}
	Log.Info(fmt.Sprintln("downloaded file", patch.GetFilepath()))
	return nil
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
	err := filepath.WalkDir(fs_path, func(path string, d fs.DirEntry, err error) error {
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
	Log.CheckError(err, "failed to walk dir")
	return myfs
}

func WriteToFile(path string, content []byte) {
	destination, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintln("os.Create:", err))
	}
	defer func() {
		err := destination.Close()
		Log.CheckError(err, "failed to close destination")
	}()

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

type LauncherConfigData struct {
	PatchMap map[PatchHash]string
	Patches  []PatchHash
}

func ReadLauncherConfig() (LauncherConfigData, []string) {
	result := LauncherConfigData{
		PatchMap: make(map[PatchHash]string),
	}
	body, err := os.ReadFile("launcherconfig.xml")
	if err != nil {
		panic(fmt.Sprintln("failed to read launcherconfig", err))
	}

	var Page BadassRoot
	err = xml.Unmarshal(body, &Page)
	Log.CheckError(err, "failed to unarmshal xml for launcherconfig.xml")

	for _, patch := range Page.PatchHistory.Patch {
		result.PatchMap[PatchHash(patch)] = ""
		result.Patches = append(result.Patches, PatchHash(patch))
	}

	str := string(body)
	file_lines := strings.Split(str, "\n")
	return result, file_lines
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

const DiscoveryUrl = "https://patch.discoverygc.com/"
const DiscoveryUrlCache = "https://disco-api.dd84ai.com/"

func RunAutopatcher(workdir string, use_cache bool) error {
	err := os.Chdir(workdir)
	Log.CheckError(err, "failed to change working directory")
	println(os.Getwd())

	discovery_url := DiscoveryUrl
	if use_cache {
		discovery_url = DiscoveryUrlCache
	}

	discovery_path_url := discovery_url + "patchlist.xml"
	resp, err := Request(discovery_path_url)
	if err != nil {
		return err
	}
	patches := ParseForPatches(discovery_url, resp.Body)

	patchhistory, file_lines := ReadLauncherConfig()

	var applied_patches []Patch

	for index, patch := range patches {
		if _, found := patchhistory.PatchMap[patch.Hash]; found {
			Log.Info(fmt.Sprintln("patch is already installed", patch))

			if index == len(patches)-1 {
				patch_marshaled, _ := json.Marshal(patch)
				err = os.WriteFile(AutopatherFilename, patch_marshaled, 0666)
				Log.CheckError(err, "failed to write patch_marshaled")
			}
			continue
		}

		err := downloadPatch(patch)
		if Log.CheckError(err, "failed to download patch") {
			return err
		}

		patch_body, _ := os.ReadFile(patch.GetFilepath())
		hash := md5.Sum(patch_body)
		md5_result := hex.EncodeToString(hash[:])
		Log.Info(fmt.Sprintln("md5_result=", md5_result))

		if md5_result != strings.ToLower(string(patch.Hash)) {
			return errors.New(fmt.Sprintln("md5 hash sum is not matching", "expected=", patch.Hash, " but bound=", md5_result))
		}

		err = Unzip(patch.GetFilepath(), patch.GetFolderPath())
		Log.CheckError(err, "failed to unzip")

		freelancer_folder := ScanCaseInsensitiveFS(".")
		patch_folder := ScanCaseInsensitiveFS(patch.GetFolderPath())

		for _, file := range patch_folder.Files {

			content, err := os.ReadFile(file.GetPath())
			if Log.CheckError(err, fmt.Sprintln("failed to read file", err)) {
				return err
			}

			relative_patch_filepath := file.GetRelPathTo(patch.GetFolderPath())

			if strings.Contains(relative_patch_filepath, ".gitignore") {
				continue
			}

			if freelancer_path, file_exists := freelancer_folder.LowerMapFiles[strings.ToLower(relative_patch_filepath)]; file_exists {
				_ = os.Remove(freelancer_path.GetPath())
			}

			relative_patch_filepath = AdjustFoldersInPath(relative_patch_filepath, freelancer_folder)

			WriteToFile(relative_patch_filepath, content)

		}

		_ = os.RemoveAll(patch.GetFolderPath())

		Log.Info("applied patch", typelog.Any("patch", patch))
		patch_marshaled, _ := json.Marshal(patch)
		err = os.WriteFile(AutopatherFilename, patch_marshaled, 0666)
		Log.CheckError(err, "failed to write file")

		applied_patches = append(applied_patches, patch)
	}

	var patch_file_start []string
	var patch_file_end []string
	for line_index, _ := range file_lines {
		if strings.Contains(file_lines[line_index], "<Patch>") && !strings.Contains(file_lines[line_index+1], "<Patch>") {
			patch_file_start = file_lines[:line_index+1]
			patch_file_end = file_lines[line_index+1:]
		}
	}
	if len(patch_file_start) == 0 {
		err_msg := "not found patch line index, where to insert"
		return errors.New(err_msg)
	}

	Log.Info("updating launcherconfig.yml")
	var new_patch_file_lines []string
	new_patch_file_lines = append(new_patch_file_lines, patch_file_start...)
	for _, patch := range applied_patches {
		new_patch_file_lines = append(new_patch_file_lines, fmt.Sprintf("    <Patch>%s</Patch>\r", patch.Hash))
	}
	new_patch_file_lines = append(new_patch_file_lines, patch_file_end...)
	err = os.WriteFile("launcherconfig.xml", []byte(strings.Join(new_patch_file_lines, "\n")), 0666)
	Log.CheckError(err, "failed updating launcherconfig.xml")
	Log.Info("finished patchdisco run")
	return nil
}

const AutopatherFilename = "autopatcher.latest_patch.json"
