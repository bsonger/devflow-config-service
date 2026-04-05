package app

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sync"
	"testing"

	"github.com/bsonger/devflow-config-service/pkg/domain"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/bsonger/devflow-service-common/loggingx"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestCreateConfigurationUsesSQLInsert(t *testing.T) {
	loggingx.Logger = zap.NewNop()

	stub := &sqlDriverStub{execResult: driver.RowsAffected(1)}
	db := openSQLStub(t, stub)
	store.InitPostgres(db)

	cfgValue := testConfiguration()
	cfg := &cfgValue
	id, err := NewConfigurationService(nil).Create(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if id != cfg.ID {
		t.Fatalf("Create returned id %s, want %s", id, cfg.ID)
	}
	if len(stub.execs) != 1 {
		t.Fatalf("exec count = %d, want 1", len(stub.execs))
	}
}

func TestListConfigurationsUsesSQLQuery(t *testing.T) {
	loggingx.Logger = zap.NewNop()

	stub := &sqlDriverStub{queryRows: &sqlRowsStub{}}
	db := openSQLStub(t, stub)
	store.InitPostgres(db)

	items, err := NewConfigurationService(nil).List(context.Background(), ConfigurationListFilter{})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("List returned %d items, want 0", len(items))
	}
	if len(stub.queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(stub.queries))
	}
}

func testConfiguration() domain.Configuration {
	id := uuid.New()
	applicationID := uuid.New()
	revisionID := uuid.New()
	return domain.Configuration{
		BaseModel:        domain.BaseModel{ID: id},
		ApplicationID:    applicationID,
		Name:             "cfg-1",
		Env:              "staging",
		SourcePath:       "applications/example-app/staging",
		LatestRevisionNo: 1,
		LatestRevisionID: &revisionID,
	}
}

type sqlDriverStub struct {
	mu         sync.Mutex
	execs      []string
	queries    []string
	execResult driver.Result
	execErr    error
	queryRows  driver.Rows
	queryErr   error
}

func (s *sqlDriverStub) Open(name string) (driver.Conn, error) {
	return &sqlConnStub{stub: s}, nil
}

type sqlConnStub struct {
	stub *sqlDriverStub
}

func (c *sqlConnStub) Prepare(query string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c *sqlConnStub) Close() error                              { return nil }
func (c *sqlConnStub) Begin() (driver.Tx, error)                 { return nil, driver.ErrSkip }

func (c *sqlConnStub) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	c.stub.mu.Lock()
	defer c.stub.mu.Unlock()
	c.stub.execs = append(c.stub.execs, query)
	return c.stub.execResult, c.stub.execErr
}

func (c *sqlConnStub) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	c.stub.mu.Lock()
	defer c.stub.mu.Unlock()
	c.stub.queries = append(c.stub.queries, query)
	if c.stub.queryRows != nil || c.stub.queryErr != nil {
		return c.stub.queryRows, c.stub.queryErr
	}
	return &sqlRowsStub{}, nil
}

type sqlRowsStub struct{}

func (r *sqlRowsStub) Columns() []string           { return nil }
func (r *sqlRowsStub) Close() error                { return nil }
func (r *sqlRowsStub) Next(_ []driver.Value) error { return io.EOF }

func openSQLStub(t *testing.T, stub *sqlDriverStub) *sql.DB {
	t.Helper()
	driverName := "config_service_sql_stub_" + uuid.NewString()
	sql.Register(driverName, stub)

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
