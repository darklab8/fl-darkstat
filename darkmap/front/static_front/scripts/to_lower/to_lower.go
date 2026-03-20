package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	workdir := "../../discovery_navmap"

	var paths []string

	// Collect all files and folders first
	err := filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == workdir {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	if err != nil {
		fmt.Println("Walk error:", err)
		return
	}

	// Sort deepest first so children rename before parents
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) > len(paths[j])
	})

	// Rename each path
	for _, oldPath := range paths {
		dir := filepath.Dir(oldPath)
		base := filepath.Base(oldPath)
		lower := strings.ToLower(base)

		if base == lower {
			continue
		}

		newPath := filepath.Join(dir, lower)

		// Skip if target already exists
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("Skip (exists): %s -> %s\n", oldPath, newPath)
			continue
		}

		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("Rename failed: %s -> %s (%v)\n", oldPath, newPath, err)
			continue
		}

		fmt.Printf("Renamed: %s -> %s\n", oldPath, newPath)
	}
}
