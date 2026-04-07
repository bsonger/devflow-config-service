package config_repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type layoutResolution struct {
	SourcePath string
	ServiceDir string
	Env        string
	Files      []string
}

func resolveLayout(rootDir, sourcePath, env string) (*layoutResolution, error) {
	normalizedSource := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(sourcePath)), "./")
	if normalizedSource == "." || normalizedSource == "" {
		return nil, ErrSourcePathNotFound
	}

	cleanEnv := strings.TrimSpace(env)
	serviceRel := normalizedSource

	if strings.HasSuffix(normalizedSource, ".yaml") {
		dir := filepath.ToSlash(filepath.Dir(normalizedSource))
		if filepath.Base(dir) == "environments" {
			serviceRel = filepath.ToSlash(filepath.Dir(dir))
			cleanEnv = strings.TrimSuffix(filepath.Base(normalizedSource), filepath.Ext(normalizedSource))
		} else {
			serviceRel = dir
		}
	} else if filepath.Base(filepath.ToSlash(filepath.Dir(normalizedSource))) == "environments" {
		serviceRel = filepath.ToSlash(filepath.Dir(filepath.Dir(normalizedSource)))
		cleanEnv = filepath.Base(normalizedSource)
	}

	if cleanEnv == "" {
		return nil, fmt.Errorf("environment is required for normalized config repo layout")
	}

	serviceDir := filepath.Join(rootDir, filepath.FromSlash(serviceRel))
	if info, err := os.Stat(serviceDir); err != nil || !info.IsDir() {
		return nil, ErrSourcePathNotFound
	}

	files := []string{
		"configuration.yaml",
		"deployment.yaml",
		"service.yaml",
		filepath.ToSlash(filepath.Join("environments", cleanEnv+".yaml")),
	}
	for _, name := range files {
		if _, err := os.Stat(filepath.Join(serviceDir, filepath.FromSlash(name))); err != nil {
			if os.IsNotExist(err) {
				return nil, ErrSourcePathNotFound
			}
			return nil, err
		}
	}

	return &layoutResolution{
		SourcePath: normalizedSource,
		ServiceDir: serviceDir,
		Env:        cleanEnv,
		Files:      files,
	}, nil
}
