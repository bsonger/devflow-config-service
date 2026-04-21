package app

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bsonger/devflow-config-service/pkg/domain"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/google/uuid"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	store.InitPostgres(db)
	cleanup := func() {
		store.InitPostgres(nil)
		db.Close()
	}
	return mock, cleanup
}

func newValidAppConfig() *domain.AppConfig {
	return &domain.AppConfig{
		ApplicationID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		EnvironmentID: "staging",
		Name:          "devflow-app-service",
		Format:        "yaml",
		Data:          "foo: bar",
	}
}

func TestCreate_WithCustomSourcePath_UsesProvidedPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.SourcePath = "custom/path/to/config"

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into configurations (
			id, application_id, name, env, description, format, data, mount_path, labels, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'[]'::jsonb,$11,$12,$13,$14,$15)
	`)).WithArgs(
		sqlmock.AnyArg(), // id
		cfg.ApplicationID,
		cfg.Name,
		cfg.EnvironmentID,
		cfg.Description,
		cfg.Format,
		cfg.Data,
		"/etc/devflow/config", // normalized mount path
		sqlmock.AnyArg(),      // labels
		"custom/path/to/config",
		sqlmock.AnyArg(), // latest_revision_no
		sqlmock.AnyArg(), // latest_revision_id
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
		sqlmock.AnyArg(), // deleted_at
	).WillReturnResult(sqlmock.NewResult(1, 1))

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestCreate_WithEmptySourcePath_DerivesPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.SourcePath = ""

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into configurations (
			id, application_id, name, env, description, format, data, mount_path, labels, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'[]'::jsonb,$11,$12,$13,$14,$15)
	`)).WithArgs(
		sqlmock.AnyArg(),
		cfg.ApplicationID,
		cfg.Name,
		cfg.EnvironmentID,
		cfg.Description,
		cfg.Format,
		cfg.Data,
		"/etc/devflow/config",
		sqlmock.AnyArg(),
		"applications/devflow-platform/services/devflow-app-service",
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestCreate_WithWhitespaceSourcePath_DerivesPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.SourcePath = "   "

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into configurations (
			id, application_id, name, env, description, format, data, mount_path, labels, source_path, files, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'[]'::jsonb,$11,$12,$13,$14,$15)
	`)).WithArgs(
		sqlmock.AnyArg(),
		cfg.ApplicationID,
		cfg.Name,
		cfg.EnvironmentID,
		cfg.Description,
		cfg.Format,
		cfg.Data,
		"/etc/devflow/config",
		sqlmock.AnyArg(),
		"applications/devflow-platform/services/devflow-app-service",
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUpdate_WithCustomSourcePath_UsesProvidedPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	appID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()

	// Get query returns current record with existing source_path
	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, application_id, name, env, description, format, data, mount_path, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations where id=$1 and deleted_at is null
	`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
		"id", "application_id", "name", "env", "description", "format", "data",
		"mount_path", "labels", "source_path", "latest_revision_no", "latest_revision_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		id, appID.String(), "devflow-app-service", "staging", "", "yaml", "foo: bar",
		"/etc/devflow/config", []byte("[]"), "old/source/path", 0, nil,
		now, now, nil,
	))

	// Update query
	mock.ExpectExec(regexp.QuoteMeta(`
		update configurations
		set application_id=$2, name=$3, env=$4, description=$5, format=$6, data=$7, mount_path=$8, labels=$9, source_path=$10, updated_at=$11
		where id=$1 and deleted_at is null
	`)).WithArgs(
		id,
		appID,
		"devflow-app-service",
		"staging",
		"",
		"yaml",
		"foo: bar",
		"/etc/devflow/config",
		sqlmock.AnyArg(), // labels
		"new/custom/path",
		sqlmock.AnyArg(), // updated_at
	).WillReturnResult(sqlmock.NewResult(1, 1))

	cfg := newValidAppConfig()
	cfg.ID = id
	cfg.SourcePath = "new/custom/path"

	svc := NewAppConfigService(nil)
	err := svc.Update(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUpdate_WithEmptySourcePath_PreservesExistingPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	appID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, application_id, name, env, description, format, data, mount_path, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations where id=$1 and deleted_at is null
	`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
		"id", "application_id", "name", "env", "description", "format", "data",
		"mount_path", "labels", "source_path", "latest_revision_no", "latest_revision_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		id, appID.String(), "devflow-app-service", "staging", "", "yaml", "foo: bar",
		"/etc/devflow/config", []byte("[]"), "existing/source/path", 0, nil,
		now, now, nil,
	))

	mock.ExpectExec(regexp.QuoteMeta(`
		update configurations
		set application_id=$2, name=$3, env=$4, description=$5, format=$6, data=$7, mount_path=$8, labels=$9, source_path=$10, updated_at=$11
		where id=$1 and deleted_at is null
	`)).WithArgs(
		id,
		appID,
		"devflow-app-service",
		"staging",
		"",
		"yaml",
		"foo: bar",
		"/etc/devflow/config",
		sqlmock.AnyArg(),
		"existing/source/path", // preserved
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	cfg := newValidAppConfig()
	cfg.ID = id
	cfg.SourcePath = ""

	svc := NewAppConfigService(nil)
	err := svc.Update(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUpdate_WithWhitespaceSourcePath_PreservesExistingPath(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	appID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, application_id, name, env, description, format, data, mount_path, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations where id=$1 and deleted_at is null
	`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
		"id", "application_id", "name", "env", "description", "format", "data",
		"mount_path", "labels", "source_path", "latest_revision_no", "latest_revision_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		id, appID.String(), "devflow-app-service", "staging", "", "yaml", "foo: bar",
		"/etc/devflow/config", []byte("[]"), "existing/source/path", 0, nil,
		now, now, nil,
	))

	mock.ExpectExec(regexp.QuoteMeta(`
		update configurations
		set application_id=$2, name=$3, env=$4, description=$5, format=$6, data=$7, mount_path=$8, labels=$9, source_path=$10, updated_at=$11
		where id=$1 and deleted_at is null
	`)).WithArgs(
		id,
		appID,
		"devflow-app-service",
		"staging",
		"",
		"yaml",
		"foo: bar",
		"/etc/devflow/config",
		sqlmock.AnyArg(),
		"existing/source/path", // preserved
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	cfg := newValidAppConfig()
	cfg.ID = id
	cfg.SourcePath = "   "

	svc := NewAppConfigService(nil)
	err := svc.Update(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUpdate_NotFound_ReturnsError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, application_id, name, env, description, format, data, mount_path, labels, source_path, latest_revision_no, latest_revision_id, created_at, updated_at, deleted_at
		from configurations where id=$1 and deleted_at is null
	`)).WithArgs(id).WillReturnError(sql.ErrNoRows)

	cfg := newValidAppConfig()
	cfg.ID = id
	cfg.SourcePath = "some/path"

	svc := NewAppConfigService(nil)
	err := svc.Update(context.Background(), cfg)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestCreate_ValidationError_MissingName(t *testing.T) {
	_, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.Name = ""

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCreate_ValidationError_MissingApplicationID(t *testing.T) {
	_, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.ApplicationID = uuid.Nil

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for missing application_id")
	}
}

func TestCreate_ValidationError_MissingEnvironmentID(t *testing.T) {
	_, cleanup := setupMockDB(t)
	defer cleanup()

	cfg := newValidAppConfig()
	cfg.EnvironmentID = ""

	svc := NewAppConfigService(nil)
	_, err := svc.Create(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for missing environment_id")
	}
}

func TestDeriveAppConfigSourcePath(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"devflow-app-service", "applications/devflow-platform/services/devflow-app-service"},
		{"my-service", "applications/devflow-platform/services/my-service"},
		{"", ""},
		{"  ", ""},
	}
	for _, tc := range cases {
		got := deriveAppConfigSourcePath(tc.name)
		if got != tc.want {
			t.Errorf("deriveAppConfigSourcePath(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestNormalizeAppConfigMountPath(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{"", "/etc/devflow/config"},
		{"  ", "/etc/devflow/config"},
		{"/custom/mount", "/custom/mount"},
		{"/custom/mount  ", "/custom/mount"},
	}
	for _, tc := range cases {
		got := normalizeAppConfigMountPath(tc.value)
		if got != tc.want {
			t.Errorf("normalizeAppConfigMountPath(%q) = %q, want %q", tc.value, got, tc.want)
		}
	}
}
