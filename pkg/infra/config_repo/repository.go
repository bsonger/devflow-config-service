package config_repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bsonger/devflow-config-service/pkg/domain"
)

var ErrSourcePathNotFound = errors.New("config repo source path not found")

var DefaultRepository *Repository

const (
	FixedRepositoryURL = "git@github.com:bsonger/devflow-config-service.git"
	FixedBranch        = "main"
)

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

func (r *Repository) ReadSnapshot(_ context.Context, sourcePath, env string) (*Snapshot, error) {
	resolved, err := resolveLayout(r.rootDir, sourcePath, env)
	if err != nil {
		return nil, err
	}

	files := make([]domain.File, 0, len(resolved.Files))
	hash := sha256.New()
	for _, name := range resolved.Files {
		content, err := os.ReadFile(filepath.Join(resolved.Dir, filepath.FromSlash(name)))
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
		SourcePath:   strings.TrimPrefix(filepath.ToSlash(resolved.SourcePath), "./"),
		SourceCommit: r.defaultRefOrMain(),
		SourceDigest: hex.EncodeToString(hash.Sum(nil)),
		Files:        files,
	}, nil
}

func (r *Repository) defaultRefOrMain() string {
	if r == nil || r.defaultRef == "" {
		return FixedBranch
	}
	return r.defaultRef
}
