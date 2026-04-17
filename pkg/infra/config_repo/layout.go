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
	files := make([]string, 0, len(entries)+1)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}
	if envFile := resolveEnvironmentOverlay(dir, env); envFile != "" {
		files = append(files, envFile)
	}
	if len(files) == 0 {
		return nil, ErrSourcePathNotFound
	}
	sort.Slice(files, func(i, j int) bool {
		leftWeight, rightWeight := layoutFileOrder(files[i]), layoutFileOrder(files[j])
		if leftWeight != rightWeight {
			return leftWeight < rightWeight
		}
		return files[i] < files[j]
	})

	return &layoutResolution{
		SourcePath: normalizedSource,
		Dir:        dir,
		Files:      files,
	}, nil
}

func resolveEnvironmentOverlay(dir, env string) string {
	trimmedEnv := strings.TrimSpace(env)
	if trimmedEnv == "" || strings.EqualFold(trimmedEnv, "base") {
		return ""
	}
	relative := filepath.ToSlash(filepath.Join("environments", trimmedEnv+".yaml"))
	if info, err := os.Stat(filepath.Join(dir, filepath.FromSlash(relative))); err == nil && !info.IsDir() {
		return relative
	}
	return ""
}

func layoutFileOrder(name string) int {
	switch filepath.ToSlash(name) {
	case "configuration.yaml":
		return 0
	case "deployment.yaml":
		return 1
	case "service.yaml":
		return 2
	default:
		if strings.HasPrefix(filepath.ToSlash(name), "environments/") {
			return 3
		}
		return 10
	}
}
