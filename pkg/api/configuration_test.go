package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsonger/devflow-config-service/pkg/app"
	"github.com/bsonger/devflow-config-service/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubConfigurationService struct {
	createFn func(context.Context, *domain.Configuration) (uuid.UUID, error)
	getFn    func(context.Context, uuid.UUID) (*domain.Configuration, error)
	updateFn func(context.Context, *domain.Configuration) error
	deleteFn func(context.Context, uuid.UUID) error
	listFn   func(context.Context, app.ConfigurationListFilter) ([]domain.Configuration, error)
	syncFn   func(context.Context, uuid.UUID) (*app.SyncResult, error)
}

func (s stubConfigurationService) Create(ctx context.Context, cfg *domain.Configuration) (uuid.UUID, error) {
	return s.createFn(ctx, cfg)
}

func (s stubConfigurationService) Get(ctx context.Context, id uuid.UUID) (*domain.Configuration, error) {
	return s.getFn(ctx, id)
}

func (s stubConfigurationService) Update(ctx context.Context, cfg *domain.Configuration) error {
	return s.updateFn(ctx, cfg)
}

func (s stubConfigurationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteFn(ctx, id)
}

func (s stubConfigurationService) List(ctx context.Context, filter app.ConfigurationListFilter) ([]domain.Configuration, error) {
	return s.listFn(ctx, filter)
}

func (s stubConfigurationService) Sync(ctx context.Context, id uuid.UUID) (*app.SyncResult, error) {
	return s.syncFn(ctx, id)
}

func TestCreateConfigurationReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{
			createFn: func(_ context.Context, cfg *domain.Configuration) (uuid.UUID, error) {
				return cfg.GetID(), nil
			},
		},
	}

	r := gin.New()
	r.POST("/api/v1/configurations", handler.Create)

	body := bytes.NewBufferString(`{"application_id":"11111111-1111-1111-1111-111111111111","name":"cfg","env":"staging","source_path":"applications/example-app/staging"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/configurations", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("got %d want %d", rec.Code, http.StatusCreated)
	}

	var payload struct {
		Data domain.Configuration `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Data.Name != "cfg" || payload.Data.Env != "staging" {
		t.Fatalf("unexpected payload: %#v", payload.Data)
	}
	if payload.Data.SourcePath != "applications/example-app/staging" {
		t.Fatalf("source_path = %q", payload.Data.SourcePath)
	}
}

func TestListConfigurationsReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{
			listFn: func(_ context.Context, filter app.ConfigurationListFilter) ([]domain.Configuration, error) {
				if filter.Name != "" {
					t.Fatalf("unexpected name filter: %q", filter.Name)
				}
				return []domain.Configuration{
					{Name: "cfg-1", Env: "staging"},
					{Name: "cfg-2", Env: "prod"},
				}, nil
			},
		},
	}

	r := gin.New()
	r.GET("/api/v1/configurations", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/configurations?page=1&page_size=20", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data       []domain.Configuration `json:"data"`
		Pagination struct {
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if len(payload.Data) != 2 {
		t.Fatalf("data len = %d, want 2", len(payload.Data))
	}
	if payload.Pagination.Page != 1 || payload.Pagination.PageSize != 20 || payload.Pagination.Total != 2 {
		t.Fatalf("unexpected pagination: %#v", payload.Pagination)
	}
}

func TestDeleteConfigurationReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{
			deleteFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		},
	}

	r := gin.New()
	r.DELETE("/api/v1/configurations/:id", handler.Delete)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/configurations/"+id.String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestGetConfigurationInvalidIDReturnsErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{},
	}

	r := gin.New()
	r.GET("/api/v1/configurations/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/configurations/not-a-uuid", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d want %d", rec.Code, http.StatusBadRequest)
	}

	var payload struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Error.Code != "invalid_argument" {
		t.Fatalf("code = %q, want %q", payload.Error.Code, "invalid_argument")
	}
}

func TestUpdateConfigurationNotFoundReturnsErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{
			updateFn: func(_ context.Context, _ *domain.Configuration) error {
				return sql.ErrNoRows
			},
		},
	}

	r := gin.New()
	r.PUT("/api/v1/configurations/:id", handler.Update)

	id := uuid.New()
	body := bytes.NewBufferString(`{"application_id":"11111111-1111-1111-1111-111111111111","name":"cfg","env":"staging","source_path":"applications/example-app/staging","latest_revision_no":1}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/configurations/"+id.String(), body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestSyncConfigurationReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	revisionID := uuid.New()
	cfgID := uuid.New()
	handler := &ConfigurationHandler{
		svc: stubConfigurationService{
			syncFn: func(_ context.Context, id uuid.UUID) (*app.SyncResult, error) {
				if id != cfgID {
					t.Fatalf("got id %s want %s", id, cfgID)
				}
				return &app.SyncResult{
					Revision: &domain.ConfigurationRevision{
						ID:              revisionID,
						ConfigurationID: cfgID,
						RevisionNo:      2,
						ContentHash:     "hash",
						SourceCommit:    "main",
						SourceDigest:    "digest",
					},
					Created: true,
				}, nil
			},
		},
	}

	r := gin.New()
	r.POST("/api/v1/configurations/:id/sync", handler.Sync)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/configurations/"+cfgID.String()+"/sync", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data domain.ConfigurationRevision `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Data.SourceCommit != "main" || payload.Data.ID != revisionID {
		t.Fatalf("unexpected revision: %#v", payload.Data)
	}
}
