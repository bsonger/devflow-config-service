package config_repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bsonger/devflow-config-service/pkg/domain"
)

var ErrSourcePathNotFound = errors.New("config repo source path not found")

var DefaultRepository *Repository

type Options struct {
	RootDir    string
	DefaultRef string
}

type Snapshot struct {
	SourcePath   string
	SourceCommit string
	SourceDigest string
	Files        []domain.File
}

type Repository struct {
	rootDir    string
	defaultRef string
}

func NewRepository(opts Options) *Repository {
	return &Repository{
		rootDir:    opts.RootDir,
		defaultRef: opts.DefaultRef,
	}
}

func (r *Repository) ReadSnapshot(_ context.Context, sourcePath string) (*Snapshot, error) {
	filesDir := filepath.Join(r.rootDir, sourcePath, "files")
	entries, err := os.ReadDir(filesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrSourcePathNotFound
		}
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	files := make([]domain.File, 0, len(names))
	hash := sha256.New()
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(filesDir, name))
		if err != nil {
			return nil, err
		}
		files = append(files, domain.File{
			Name:    name,
			Content: string(content),
		})
		hash.Write([]byte(name))
		hash.Write([]byte{'\n'})
		hash.Write(content)
		hash.Write([]byte{'\n'})
	}

	return &Snapshot{
		SourcePath:   strings.TrimPrefix(filepath.ToSlash(sourcePath), "./"),
		SourceCommit: r.defaultRefOrMain(),
		SourceDigest: hex.EncodeToString(hash.Sum(nil)),
		Files:        files,
	}, nil
}

func (r *Repository) defaultRefOrMain() string {
	if r == nil || r.defaultRef == "" {
		return "main"
	}
	return r.defaultRef
}
