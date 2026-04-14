package config_repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bsonger/devflow-config-service/pkg/domain"
)

var ErrSourcePathNotFound = errors.New("config repo source path not found")
var ErrRepositorySyncFailed = errors.New("config repo sync failed")

var DefaultRepository *Repository

const (
	FixedRepositoryURL = "git@github.com:bsonger/devflow-config-repo.git"
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

type gitSyncer interface {
	Sync(ctx context.Context, rootDir, ref string) (string, error)
}

type Repository struct {
	rootDir    string
	defaultRef string
	syncer     gitSyncer
}

func NewRepository(opts Options) *Repository {
	return &Repository{
		rootDir:    opts.RootDir,
		defaultRef: opts.DefaultRef,
		syncer:     commandGitSyncer{},
	}
}

func (r *Repository) ReadSnapshot(ctx context.Context, sourcePath, env string) (*Snapshot, error) {
	sourceCommit := r.defaultRefOrMain()
	if commit, err := r.sync(ctx); err != nil {
		return nil, err
	} else if commit != "" {
		sourceCommit = commit
	}
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
		SourceCommit: sourceCommit,
		SourceDigest: hex.EncodeToString(hash.Sum(nil)),
		Files:        files,
	}, nil
}

func (r *Repository) sync(ctx context.Context) (string, error) {
	if r == nil || r.rootDir == "" || r.syncer == nil {
		return "", nil
	}
	if _, err := os.Stat(filepath.Join(r.rootDir, ".git")); err != nil {
		return "", nil
	}
	commit, err := r.syncer.Sync(ctx, r.rootDir, r.defaultRefOrMain())
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRepositorySyncFailed, err)
	}
	return commit, nil
}

func (r *Repository) defaultRefOrMain() string {
	if r == nil || r.defaultRef == "" {
		return FixedBranch
	}
	return r.defaultRef
}

type commandGitSyncer struct{}

func (commandGitSyncer) Sync(ctx context.Context, rootDir, ref string) (string, error) {
	pull := exec.CommandContext(ctx, "git", "-C", rootDir, "pull", "--ff-only", "origin", ref)
	pull.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := pull.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git pull origin %s: %w: %s", ref, err, strings.TrimSpace(string(output)))
	}

	head := exec.CommandContext(ctx, "git", "-C", rootDir, "rev-parse", "HEAD")
	head.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err = head.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}
