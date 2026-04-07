package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/bsonger/devflow-config-service/pkg/domain"
	configrepo "github.com/bsonger/devflow-config-service/pkg/infra/config_repo"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/google/uuid"
)

var ErrConfigSourceNotFound = errors.New("configuration source path not found")
var ErrConfigRepositoryUnavailable = errors.New("configuration repository is not configured")

type ConfigRepository interface {
	ReadSnapshot(ctx context.Context, sourcePath, env string) (*configrepo.Snapshot, error)
}

type SyncResult struct {
	Revision *domain.ConfigurationRevision
	Created  bool
}

func ConfigureConfigRepository(repo ConfigRepository) {
	ConfigurationService.repo = repo
}

func (s *configurationService) Sync(ctx context.Context, id uuid.UUID) (*SyncResult, error) {
	if s.repo == nil {
		return nil, ErrConfigRepositoryUnavailable
	}

	cfg, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	snapshot, err := s.repo.ReadSnapshot(ctx, cfg.SourcePath, cfg.Env)
	if err != nil {
		if errors.Is(err, configrepo.ErrSourcePathNotFound) {
			return nil, ErrConfigSourceNotFound
		}
		return nil, err
	}

	latest, err := s.getLatestRevision(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if latest != nil && latest.SourceDigest == snapshot.SourceDigest {
		return &SyncResult{Revision: latest, Created: false}, nil
	}

	revision := &domain.ConfigurationRevision{
		ID:              uuid.New(),
		ConfigurationID: id,
		RevisionNo:      cfg.LatestRevisionNo + 1,
		Files:           snapshot.Files,
		ContentHash:     snapshot.SourceDigest,
		SourceCommit:    snapshot.SourceCommit,
		SourceDigest:    snapshot.SourceDigest,
		CreatedAt:       time.Now(),
	}
	if latest == nil {
		revision.RevisionNo = 1
	}

	if err := s.insertRevision(ctx, revision); err != nil {
		return nil, err
	}
	if err := s.updateLatestRevision(ctx, cfg, revision); err != nil {
		return nil, err
	}

	return &SyncResult{Revision: revision, Created: true}, nil
}

func (s *configurationService) getLatestRevision(ctx context.Context, configurationID uuid.UUID) (*domain.ConfigurationRevision, error) {
	return scanConfigurationRevision(store.DB().QueryRowContext(ctx, `
		select id, configuration_id, revision_no, files, content_hash, source_commit, source_digest, message, created_by, created_at
		from configuration_revisions
		where configuration_id = $1
		order by revision_no desc
		limit 1
	`, configurationID))
}

func (s *configurationService) insertRevision(ctx context.Context, revision *domain.ConfigurationRevision) error {
	filesJSON, err := json.Marshal(revision.Files)
	if err != nil {
		return err
	}

	_, err = store.DB().ExecContext(ctx, `
		insert into configuration_revisions (
			id, configuration_id, revision_no, files, content_hash, source_commit, source_digest, message, created_by, created_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, revision.ID, revision.ConfigurationID, revision.RevisionNo, filesJSON, revision.ContentHash, revision.SourceCommit, revision.SourceDigest, revision.Message, revision.CreatedBy, revision.CreatedAt)
	return err
}

func (s *configurationService) updateLatestRevision(ctx context.Context, cfg *domain.Configuration, revision *domain.ConfigurationRevision) error {
	cfg.LatestRevisionNo = revision.RevisionNo
	cfg.LatestRevisionID = &revision.ID
	cfg.WithUpdateDefault()

	_, err := store.DB().ExecContext(ctx, `
		update configurations
		set latest_revision_no=$2, latest_revision_id=$3, updated_at=$4
		where id = $1 and deleted_at is null
	`, cfg.ID, cfg.LatestRevisionNo, cfg.LatestRevisionID, cfg.UpdatedAt)
	return err
}

func scanConfigurationRevision(scanner interface {
	Scan(dest ...any) error
}) (*domain.ConfigurationRevision, error) {
	var (
		revision  domain.ConfigurationRevision
		filesJSON []byte
	)

	if err := scanner.Scan(
		&revision.ID,
		&revision.ConfigurationID,
		&revision.RevisionNo,
		&filesJSON,
		&revision.ContentHash,
		&revision.SourceCommit,
		&revision.SourceDigest,
		&revision.Message,
		&revision.CreatedBy,
		&revision.CreatedAt,
	); err != nil {
		return nil, err
	}
	if len(filesJSON) > 0 {
		if err := json.Unmarshal(filesJSON, &revision.Files); err != nil {
			return nil, err
		}
	}
	return &revision, nil
}
