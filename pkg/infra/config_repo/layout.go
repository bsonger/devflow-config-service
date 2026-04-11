package config_repo

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type layoutResolution struct {
	SourcePath string
	Dir        string
	Files      []string
}

func resolveLayout(rootDir, sourcePath, env string) (*layoutResolution, error) {
	normalizedSource := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(sourcePath)), "./")
	if normalizedSource == "." || normalizedSource == "" {
		return nil, ErrSourcePathNotFound
	}
	dir := filepath.Join(rootDir, filepath.FromSlash(normalizedSource))
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil, ErrSourcePathNotFound
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}
	if len(files) == 0 {
		return nil, ErrSourcePathNotFound
	}
	sort.Strings(files)

	return &layoutResolution{
		SourcePath: normalizedSource,
		Dir:        dir,
		Files:      files,
	}, nil
}
