package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bsonger/devflow-config-service/pkg/domain"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/bsonger/devflow-service-common/loggingx"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ConfigurationService = NewConfigurationService(nil)

type ConfigurationListFilter struct {
	IncludeDeleted bool
	Name           string
}

type configurationService struct {
	repo ConfigRepository
}

func NewConfigurationService(repo ConfigRepository) *configurationService {
	return &configurationService{repo: repo}
}

func (s *configurationService) Create(ctx context.Context, cfg *domain.Configuration) (uuid.UUID, error) {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "create_configuration"),
	)

	_, err := store.DB().ExecContext(ctx, `
		insert into configurations (
			id, application_id, name, env, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, cfg.ID, nullableUUID(cfg.ApplicationID), cfg.Name, cfg.Env, cfg.SourcePath, marshalJSON(cfg.Files), cfg.LatestRevisionNo, nullableUUIDPtr(cfg.LatestRevisionID), cfg.CreatedAt, cfg.UpdatedAt, cfg.DeletedAt)
	if err != nil {
		log.Error("create configuration failed", zap.Error(err))
		return uuid.Nil, err
	}

	log.Info("configuration created", zap.String("configuration_id", cfg.GetID().String()))
	return cfg.GetID(), nil
}

func (s *configurationService) Get(ctx context.Context, id uuid.UUID) (*domain.Configuration, error) {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "get_configuration"),
		zap.String("configuration_id", id.String()),
	)

	cfg, err := scanConfiguration(store.DB().QueryRowContext(ctx, `
		select id, application_id, name, env, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations
		where id = $1 and deleted_at is null
	`, id))
	if err != nil {
		log.Error("get configuration failed", zap.Error(err))
		return nil, err
	}

	log.Debug("configuration fetched", zap.String("configuration_name", cfg.Name))
	return cfg, nil
}

func (s *configurationService) Update(ctx context.Context, cfg *domain.Configuration) error {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "update_configuration"),
		zap.String("configuration_id", cfg.GetID().String()),
	)

	current, err := s.Get(ctx, cfg.GetID())
	if err != nil {
		log.Error("load configuration failed", zap.Error(err))
		return err
	}

	cfg.CreatedAt = current.CreatedAt
	cfg.DeletedAt = current.DeletedAt
	cfg.WithUpdateDefault()

	result, err := store.DB().ExecContext(ctx, `
		update configurations
		set application_id=$2, name=$3, env=$4, source_path=$5, files=$6, latest_revision_no=$7, latest_revision_id=$8, updated_at=$9, deleted_at=$10
		where id = $1 and deleted_at is null
	`, cfg.ID, nullableUUID(cfg.ApplicationID), cfg.Name, cfg.Env, cfg.SourcePath, marshalJSON(cfg.Files), cfg.LatestRevisionNo, nullableUUIDPtr(cfg.LatestRevisionID), cfg.UpdatedAt, cfg.DeletedAt)
	if err != nil {
		log.Error("update configuration failed", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Debug("configuration updated", zap.String("configuration_name", cfg.Name))
	return nil
}

func (s *configurationService) Delete(ctx context.Context, id uuid.UUID) error {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "delete_configuration"),
		zap.String("configuration_id", id.String()),
	)

	now := time.Now()
	result, err := store.DB().ExecContext(ctx, `
		update configurations
		set deleted_at=$2, updated_at=$2
		where id = $1 and deleted_at is null
	`, id, now)
	if err != nil {
		log.Error("delete configuration failed", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Info("configuration deleted")
	return nil
}

func (s *configurationService) List(ctx context.Context, filter ConfigurationListFilter) ([]domain.Configuration, error) {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "list_configurations"),
		zap.Any("filter", filter),
	)

	query := `
		select id, application_id, name, env, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations
	`
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if !filter.IncludeDeleted {
		clauses = append(clauses, "deleted_at is null")
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
		log.Error("list configurations failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	cfgs := make([]domain.Configuration, 0)
	for rows.Next() {
		cfg, err := scanConfiguration(rows)
		if err != nil {
			return nil, err
		}
		cfgs = append(cfgs, *cfg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Debug("configurations listed", zap.Int("count", len(cfgs)))
	return cfgs, nil
}

func scanConfiguration(scanner interface {
	Scan(dest ...any) error
}) (*domain.Configuration, error) {
	var (
		cfg              domain.Configuration
		applicationID    sql.NullString
		filesJSON        []byte
		latestRevisionID sql.NullString
		deletedAt        sql.NullTime
	)

	if err := scanner.Scan(
		&cfg.ID,
		&applicationID,
		&cfg.Name,
		&cfg.Env,
		&cfg.SourcePath,
		&filesJSON,
		&cfg.LatestRevisionNo,
		&latestRevisionID,
		&cfg.CreatedAt,
		&cfg.UpdatedAt,
		&deletedAt,
	); err != nil {
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
	if len(filesJSON) > 0 {
		if err := json.Unmarshal(filesJSON, &cfg.Files); err != nil {
			return nil, err
		}
	}
	if deletedAt.Valid {
		cfg.DeletedAt = &deletedAt.Time
	}

	return &cfg, nil
}

func marshalJSON(value any) []byte {
	if value == nil {
		return []byte("[]")
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return []byte("[]")
	}
	return payload
}

func nullableUUID(id uuid.UUID) any {
	if id == uuid.Nil {
		return nil
	}
	return id
}

func nullableUUIDPtr(id *uuid.UUID) any {
	if id == nil || *id == uuid.Nil {
		return nil
	}
	return *id
}

func placeholderClause(column string, position int) string {
	return column + " = $" + strconv.Itoa(position)
}
