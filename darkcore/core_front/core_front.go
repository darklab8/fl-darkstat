package core_front

import (
	"crypto/md5"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type StaticFilesystem struct {
	Files         []core_types.StaticFile
	relPathToFile map[utils_types.FilePath]core_types.StaticFile
}

func (fs StaticFilesystem) GetFileByRelPath(rel_path utils_types.FilePath) core_types.StaticFile {
	file, ok := fs.relPathToFile[rel_path]

	if !ok {
		logus.Log.Panic("expected file found by relpath", typelog.Any("relpath", rel_path))
	}

	return file
}

func printFilesForDebug(filesystem embed.FS) {
	fs.WalkDir(filesystem, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			st, _ := fs.Stat(filesystem, p)
			r, _ := filesystem.Open(p)
			defer r.Close()

			var buf [md5.Size]byte
			n, _ := io.ReadFull(r, buf[:])

			h := md5.New()
			_, _ = io.Copy(h, r)
			s := h.Sum(nil)

			fmt.Printf("%s %d %x %x\n", p, st.Size(), buf[:n], s)
		}
		return nil
	})
}

type embedWalkState struct {
	rootFolder        utils_types.FilePath
	relFolder         utils_types.FilePath
	allowedExtensions []string
	isNotRecursive    bool
}

func walkEmbedDir(filesystem embed.FS, st embedWalkState) []utils_types.File {
	dirPath := filepath.ToSlash(st.rootFolder.ToString())
	files, err := filesystem.ReadDir(dirPath)
	var result []utils_types.File
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if st.isNotRecursive {
				continue
			}
			sub := st
			sub.relFolder = st.relFolder.Join(f.Name())
			sub.rootFolder = st.rootFolder.Join(f.Name())
			result = append(result,
				walkEmbedDir(filesystem, sub)...,
			)
		} else {
			splitted := strings.Split(f.Name(), ".")
			var extension string
			if len(splitted) > 0 {
				extension = splitted[len(splitted)-1]
			}

			path := st.rootFolder.Join(f.Name())
			requested := strings.ReplaceAll(path.ToString(), "\\", "/")
			content, err := filesystem.ReadFile(requested)
			if err != nil {
				printFilesForDebug(filesystem)
				fmt.Println(err.Error(), "failed to read file from embeded fs of",
					"path=", path,
					"requested", requested,
				)
			}

			is_allowed_extension := false
			for _, allowed_extension := range st.allowedExtensions {
				if allowed_extension == extension {
					is_allowed_extension = true
				}
			}
			if !is_allowed_extension {
				continue
			}

			relp := filepath.ToSlash(string(st.relFolder.Join(f.Name())))
			result = append(result, utils_types.File{
				Relpath:   utils_types.FilePath(relp),
				Name:      f.Name(),
				Extension: extension,
				Content:   string(content),
			})

		}

	}
	return result
}

func getFilesFromEmbed(filesystem embed.FS, params utils_types.GetFilesParams) []utils_types.File {
	st := embedWalkState{
		rootFolder:        params.RootFolder,
		relFolder:         "",
		allowedExtensions: params.AllowedExtensions,
		isNotRecursive:    params.IsNotRecursive,
	}
	if len(st.allowedExtensions) == 0 {
		st.allowedExtensions = []string{"js", "css", "png", "jpeg", "json"}
	}
	if st.rootFolder == "" {
		st.rootFolder = "."
	}
	return walkEmbedDir(filesystem, st)
}

func GetFiles(fs embed.FS, params utils_types.GetFilesParams) StaticFilesystem {
	files := getFilesFromEmbed(fs, params)
	var filesystem StaticFilesystem = StaticFilesystem{
		relPathToFile: make(map[utils_types.FilePath]core_types.StaticFile),
	}

	for _, file := range files {
		var static_file_kind core_types.StaticFileKind

		switch file.Extension {
		case "js":
			static_file_kind = core_types.StaticFileJS
		case "css":
			static_file_kind = core_types.StaticFileCSS
		case "ico":
			static_file_kind = core_types.StaticFileIco
		}

		new_file := core_types.StaticFile{
			Filename: string(file.Relpath),
			Kind:     static_file_kind,
			Content:  string(file.Content),
		}
		filesystem.Files = append(filesystem.Files, new_file)
		filesystem.relPathToFile[file.Relpath] = new_file
	}
	return filesystem
}
