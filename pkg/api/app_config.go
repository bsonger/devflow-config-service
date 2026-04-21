package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/bsonger/devflow-config-service/pkg/app"
	"github.com/bsonger/devflow-config-service/pkg/domain"
	"github.com/bsonger/devflow-service-common/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var AppConfigRouteApi = NewAppConfigHandler()

type appConfigService interface {
	Create(ctx context.Context, cfg *domain.AppConfig) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.AppConfig, error)
	Update(ctx context.Context, cfg *domain.AppConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter app.AppConfigListFilter) ([]domain.AppConfig, error)
	Sync(ctx context.Context, id uuid.UUID) (*app.AppConfigSyncResult, error)
}

type AppConfigHandler struct{ svc appConfigService }

func NewAppConfigHandler() *AppConfigHandler { return &AppConfigHandler{svc: app.AppConfigService} }

// @Summary Create app config
// @Tags AppConfig
// @Accept json
// @Produce json
// @Param data body domain.AppConfigInput true "AppConfig data"
// @Success 201 {object} httpx.DataResponse[domain.AppConfig]
// @Router /api/v1/app-configs [post]
func (h *AppConfigHandler) Create(c *gin.Context) {
	var req domain.AppConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	item := &domain.AppConfig{
		ApplicationID: req.ApplicationID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Description:   req.Description,
		Format:        req.Format,
		Data:          req.Data,
		MountPath:     req.MountPath,
		Labels:        req.Labels,
		SourcePath:    req.SourcePath,
	}
	item.WithCreateDefault()
	if _, err := h.svc.Create(c.Request.Context(), item); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	httpx.WriteData(c, http.StatusCreated, item)
}

// @Summary Get app config
// @Tags AppConfig
// @Produce json
// @Param id path string true "AppConfig ID"
// @Success 200 {object} httpx.DataResponse[domain.AppConfig]
// @Router /api/v1/app-configs/{id} [get]
func (h *AppConfigHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}
	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}
	httpx.WriteData(c, http.StatusOK, item)
}

// @Summary Update app config
// @Tags AppConfig
// @Accept json
// @Param id path string true "AppConfig ID"
// @Param data body domain.AppConfigInput true "AppConfig data"
// @Success 204
// @Router /api/v1/app-configs/{id} [put]
func (h *AppConfigHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}
	var req domain.AppConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	item := &domain.AppConfig{
		ApplicationID: req.ApplicationID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Description:   req.Description,
		Format:        req.Format,
		Data:          req.Data,
		MountPath:     req.MountPath,
		Labels:        req.Labels,
		SourcePath:    req.SourcePath,
	}
	item.SetID(id)
	if err := h.svc.Update(c.Request.Context(), item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	httpx.WriteNoContent(c)
}

// @Summary Delete app config
// @Tags AppConfig
// @Param id path string true "AppConfig ID"
// @Success 204
// @Router /api/v1/app-configs/{id} [delete]
func (h *AppConfigHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}
	httpx.WriteNoContent(c)
}

// @Summary List app configs
// @Tags AppConfig
// @Produce json
// @Param application_id query string false "Application ID"
// @Param environment_id query string false "Environment ID"
// @Param name query string false "Name"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpx.ListResponse[domain.AppConfig]
// @Router /api/v1/app-configs [get]
func (h *AppConfigHandler) List(c *gin.Context) {
	var filter app.AppConfigListFilter
	if appID := c.Query("application_id"); appID != "" {
		id, err := uuid.Parse(appID)
		if err != nil {
			httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid application_id", nil)
			return
		}
		filter.ApplicationID = &id
	}
	filter.EnvironmentID = c.Query("environment_id")
	filter.Name = c.Query("name")
	filter.IncludeDeleted = httpx.IncludeDeleted(c)
	items, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}
	paging, err := httpx.ParsePagination(c)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	total := len(items)
	items = httpx.PaginateSlice(items, paging)
	httpx.WriteList(c, http.StatusOK, items, paging, total)
}

// @Summary Sync app config from fixed config repo
// @Tags AppConfig
// @Produce json
// @Param id path string true "AppConfig ID"
// @Success 200 {object} httpx.DataResponse[domain.AppConfigRevision]
// @Router /api/v1/app-configs/{id}/sync-from-repo [post]
func (h *AppConfigHandler) Sync(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}
	result, err := h.svc.Sync(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
		case errors.Is(err, app.ErrConfigSourceNotFound), errors.Is(err, app.ErrConfigRepositoryUnavailable), errors.Is(err, app.ErrConfigRepositorySyncFailed):
			httpx.WriteError(c, http.StatusFailedDependency, "failed_precondition", err.Error(), nil)
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		}
		return
	}
	httpx.WriteData(c, http.StatusOK, result.Revision)
}
