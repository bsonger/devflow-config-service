package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bsonger/devflow-config-service/pkg/domain"
	configrepo "github.com/bsonger/devflow-config-service/pkg/infra/config_repo"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/google/uuid"
)

var AppConfigService = NewAppConfigService(nil)

var ErrConfigSourceNotFound = errors.New("configuration source path not found")
var ErrConfigRepositoryUnavailable = errors.New("configuration repository is not configured")
var ErrConfigRepositorySyncFailed = errors.New("configuration repository sync failed")

type AppConfigListFilter struct {
	ApplicationID  *uuid.UUID
	EnvironmentID  string
	IncludeDeleted bool
	Name           string
}

type AppConfigSyncResult struct {
	Revision *domain.AppConfigRevision
	Created  bool
}

type appConfigRepository interface {
	ReadSnapshot(ctx context.Context, sourcePath, env string) (*configrepo.Snapshot, error)
}

type appConfigService struct {
	repo appConfigRepository
}

func NewAppConfigService(repo appConfigRepository) *appConfigService {
	return &appConfigService{repo: repo}
}

func ConfigureAppConfigRepository(repo appConfigRepository) {
	AppConfigService.repo = repo
}

func (s *appConfigService) Create(ctx context.Context, cfg *domain.AppConfig) (uuid.UUID, error) {
	if err := validateAppConfig(cfg); err != nil {
		return uuid.Nil, err
	}
	cfg.SourcePath = deriveAppConfigSourcePath(cfg.Name)
	_, err := store.DB().ExecContext(ctx, `
		insert into configurations (
			id, application_id, name, env, description, format, data, labels, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,'[]'::jsonb,$10,$11,$12,$13,$14)
	`, cfg.ID, cfg.ApplicationID, cfg.Name, cfg.EnvironmentID, cfg.Description, cfg.Format, cfg.Data, marshalJSON(cfg.Labels), cfg.SourcePath, cfg.LatestRevisionNo, nullableUUIDPtr(cfg.LatestRevisionID), cfg.CreatedAt, cfg.UpdatedAt, cfg.DeletedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return cfg.ID, nil
}

func (s *appConfigService) Get(ctx context.Context, id uuid.UUID) (*domain.AppConfig, error) {
	cfg, err := scanAppConfig(store.DB().QueryRowContext(ctx, `
		select id, application_id, name, env, description, format, data, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations where id=$1 and deleted_at is null
	`, id))
	if err != nil {
		return nil, err
	}
	if cfg.LatestRevisionID == nil || *cfg.LatestRevisionID == uuid.Nil {
		return cfg, nil
	}
	revision, err := s.getRevision(ctx, *cfg.LatestRevisionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cfg, nil
		}
		return nil, err
	}
	cfg.Files = revision.Files
	cfg.RenderedConfigMap = revision.RenderedConfigMap
	cfg.SourceCommit = revision.SourceCommit
	return cfg, nil
}

func (s *appConfigService) Update(ctx context.Context, cfg *domain.AppConfig) error {
	if err := validateAppConfig(cfg); err != nil {
		return err
	}
	current, err := s.Get(ctx, cfg.ID)
	if err != nil {
		return err
	}
	cfg.CreatedAt = current.CreatedAt
	cfg.DeletedAt = current.DeletedAt
	cfg.SourcePath = deriveAppConfigSourcePath(cfg.Name)
	cfg.WithUpdateDefault()
	result, err := store.DB().ExecContext(ctx, `
		update configurations
		set application_id=$2, name=$3, env=$4, description=$5, format=$6, data=$7, labels=$8, source_path=$9, updated_at=$10
		where id=$1 and deleted_at is null
	`, cfg.ID, cfg.ApplicationID, cfg.Name, cfg.EnvironmentID, cfg.Description, cfg.Format, cfg.Data, marshalJSON(cfg.Labels), cfg.SourcePath, cfg.UpdatedAt)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *appConfigService) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result, err := store.DB().ExecContext(ctx, `
		update configurations set deleted_at=$2, updated_at=$2
		where id=$1 and deleted_at is null
	`, id, now)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result)
}

func (s *appConfigService) List(ctx context.Context, filter AppConfigListFilter) ([]domain.AppConfig, error) {
	query := `
		select id, application_id, name, env, description, format, data, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations
	`
	clauses := make([]string, 0, 4)
	args := make([]any, 0, 4)
	if !filter.IncludeDeleted {
		clauses = append(clauses, "deleted_at is null")
	}
	if filter.ApplicationID != nil {
		args = append(args, *filter.ApplicationID)
		clauses = append(clauses, placeholderClause("application_id", len(args)))
	}
	if filter.EnvironmentID != "" {
		args = append(args, filter.EnvironmentID)
		clauses = append(clauses, placeholderClause("env", len(args)))
	}
	if filter.Name != "" {
		args = append(args, filter.Name)
		clauses = append(clauses, placeholderClause("name", len(args)))
	}
	if len(clauses) > 0 {
		query += " where " + strings.Join(clauses, " and ")
	}
	query += " order by created_at desc"
	rows, err := store.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.AppConfig, 0)
	for rows.Next() {
		item, err := scanAppConfig(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

func (s *appConfigService) Sync(ctx context.Context, id uuid.UUID) (*AppConfigSyncResult, error) {
	if s.repo == nil {
		return nil, ErrConfigRepositoryUnavailable
	}
	cfg, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	snapshot, err := s.repo.ReadSnapshot(ctx, cfg.SourcePath, cfg.EnvironmentID)
	if err != nil {
		if errors.Is(err, configrepo.ErrSourcePathNotFound) {
			fallbackPath := deriveAppConfigSourcePath(cfg.Name)
			if fallbackPath != "" && fallbackPath != cfg.SourcePath {
				snapshot, err = s.repo.ReadSnapshot(ctx, fallbackPath, cfg.EnvironmentID)
				if err == nil {
					cfg.SourcePath = fallbackPath
					if updateErr := s.updateSourcePath(ctx, cfg.ID, fallbackPath); updateErr != nil {
						return nil, updateErr
					}
				}
			}
		}
		if err != nil {
			if errors.Is(err, configrepo.ErrSourcePathNotFound) {
				return nil, ErrConfigSourceNotFound
			}
			if errors.Is(err, configrepo.ErrRepositorySyncFailed) {
				return nil, fmt.Errorf("%w: %v", ErrConfigRepositorySyncFailed, err)
			}
			return nil, err
		}
	}
	latest, err := s.getLatestRevision(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if latest != nil && latest.SourceDigest == snapshot.SourceDigest {
		return &AppConfigSyncResult{Revision: latest, Created: false}, nil
	}
	revisionNo := 1
	if latest != nil {
		revisionNo = latest.RevisionNo + 1
	}
	revision := &domain.AppConfigRevision{
		ID:                uuid.New(),
		AppConfigID:       id,
		RevisionNo:        revisionNo,
		Files:             snapshot.Files,
		RenderedConfigMap: renderConfigMap(snapshot.Files),
		ContentHash:       snapshot.SourceDigest,
		SourceCommit:      snapshot.SourceCommit,
		SourceDigest:      snapshot.SourceDigest,
		CreatedAt:         time.Now().Format(time.RFC3339),
	}
	if err := s.insertRevision(ctx, revision); err != nil {
		return nil, err
	}
	if err := s.updateLatestRevision(ctx, cfg, revision); err != nil {
		return nil, err
	}
	return &AppConfigSyncResult{Revision: revision, Created: true}, nil
}

func (s *appConfigService) getLatestRevision(ctx context.Context, appConfigID uuid.UUID) (*domain.AppConfigRevision, error) {
	return scanAppConfigRevision(store.DB().QueryRowContext(ctx, `
		select id, configuration_id, revision_no, files, rendered_configmap, content_hash, source_commit, source_digest, created_at
		from configuration_revisions
		where configuration_id=$1
		order by revision_no desc limit 1
	`, appConfigID))
}

func (s *appConfigService) getRevision(ctx context.Context, id uuid.UUID) (*domain.AppConfigRevision, error) {
	return scanAppConfigRevision(store.DB().QueryRowContext(ctx, `
		select id, configuration_id, revision_no, files, rendered_configmap, content_hash, source_commit, source_digest, created_at
		from configuration_revisions
		where id=$1
	`, id))
}

func (s *appConfigService) insertRevision(ctx context.Context, revision *domain.AppConfigRevision) error {
	filesJSON, err := json.Marshal(revision.Files)
	if err != nil {
		return err
	}
	renderedJSON, err := json.Marshal(revision.RenderedConfigMap)
	if err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, revision.CreatedAt)
	if err != nil {
		return err
	}
	_, err = store.DB().ExecContext(ctx, `
		insert into configuration_revisions (
			id, configuration_id, revision_no, files, rendered_configmap, content_hash, source_commit, source_digest, message, created_by, created_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,'','',$9)
	`, revision.ID, revision.AppConfigID, revision.RevisionNo, filesJSON, renderedJSON, revision.ContentHash, revision.SourceCommit, revision.SourceDigest, createdAt)
	return err
}

func (s *appConfigService) updateLatestRevision(ctx context.Context, cfg *domain.AppConfig, revision *domain.AppConfigRevision) error {
	cfg.LatestRevisionNo = revision.RevisionNo
	cfg.LatestRevisionID = &revision.ID
	cfg.WithUpdateDefault()
	_, err := store.DB().ExecContext(ctx, `
		update configurations
		set latest_revision_no=$2, latest_revision_id=$3, updated_at=$4
		where id=$1 and deleted_at is null
	`, cfg.ID, cfg.LatestRevisionNo, cfg.LatestRevisionID, cfg.UpdatedAt)
	return err
}

func (s *appConfigService) updateSourcePath(ctx context.Context, id uuid.UUID, sourcePath string) error {
	_, err := store.DB().ExecContext(ctx, `
		update configurations
		set source_path=$2, updated_at=$3
		where id=$1 and deleted_at is null
	`, id, sourcePath, time.Now())
	return err
}

func validateAppConfig(cfg *domain.AppConfig) error {
	if cfg == nil {
		return errors.New("app_config is required")
	}
	if errs := validateAppConfigInput(cfg.ApplicationID, cfg.EnvironmentID); len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	if strings.TrimSpace(cfg.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

func validateAppConfigInput(applicationID uuid.UUID, environmentID string) []string {
	var errs []string
	if applicationID == uuid.Nil {
		errs = append(errs, "application_id is required")
	}
	if strings.TrimSpace(environmentID) == "" {
		errs = append(errs, "environment_id is required")
	}
	return errs
}

func deriveAppConfigSourcePath(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	return fmt.Sprintf("applications/devflow-platform/services/%s", trimmed)
}

func renderConfigMap(files []domain.File) domain.RenderedConfigMap {
	data := make(map[string]string, len(files))
	for _, file := range files {
		data[file.Name] = file.Content
	}
	return domain.RenderedConfigMap{Data: data}
}

func scanAppConfig(scanner interface{ Scan(dest ...any) error }) (*domain.AppConfig, error) {
	var (
		cfg              domain.AppConfig
		applicationID    sql.NullString
		labelsJSON       []byte
		latestRevisionID sql.NullString
		deletedAt        sql.NullTime
	)
	if err := scanner.Scan(&cfg.ID, &applicationID, &cfg.Name, &cfg.EnvironmentID, &cfg.Description, &cfg.Format, &cfg.Data, &labelsJSON, &cfg.SourcePath, &cfg.LatestRevisionNo, &latestRevisionID, &cfg.CreatedAt, &cfg.UpdatedAt, &deletedAt); err != nil {
		return nil, err
	}
	if applicationID.Valid {
		parsed, err := uuid.Parse(applicationID.String)
		if err != nil {
			return nil, err
		}
		cfg.ApplicationID = parsed
	}
	if latestRevisionID.Valid {
		parsed, err := uuid.Parse(latestRevisionID.String)
		if err != nil {
			return nil, err
		}
		cfg.LatestRevisionID = &parsed
	}
	if len(labelsJSON) > 0 {
		_ = json.Unmarshal(labelsJSON, &cfg.Labels)
	}
	if deletedAt.Valid {
		cfg.DeletedAt = &deletedAt.Time
	}
	return &cfg, nil
}

func scanAppConfigRevision(scanner interface{ Scan(dest ...any) error }) (*domain.AppConfigRevision, error) {
	var (
		revision     domain.AppConfigRevision
		filesJSON    []byte
		renderedJSON []byte
		createdAt    time.Time
	)
	if err := scanner.Scan(&revision.ID, &revision.AppConfigID, &revision.RevisionNo, &filesJSON, &renderedJSON, &revision.ContentHash, &revision.SourceCommit, &revision.SourceDigest, &createdAt); err != nil {
		return nil, err
	}
	if len(filesJSON) > 0 {
		if err := json.Unmarshal(filesJSON, &revision.Files); err != nil {
			return nil, err
		}
	}
	if len(renderedJSON) > 0 {
		if err := json.Unmarshal(renderedJSON, &revision.RenderedConfigMap); err != nil {
			return nil, err
		}
	}
	revision.CreatedAt = createdAt.Format(time.RFC3339)
	return &revision, nil
}

func ensureRowsAffected(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
