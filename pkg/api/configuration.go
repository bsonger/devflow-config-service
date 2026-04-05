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

var ConfigurationRouteApi = NewConfigurationHandler()

type configurationService interface {
	Create(ctx context.Context, cfg *domain.Configuration) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Configuration, error)
	Update(ctx context.Context, cfg *domain.Configuration) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter app.ConfigurationListFilter) ([]domain.Configuration, error)
	Sync(ctx context.Context, id uuid.UUID) (*app.SyncResult, error)
}

type ConfigurationHandler struct {
	svc configurationService
}

type CreateConfigurationRequest struct {
	ApplicationID uuid.UUID `json:"application_id"`
	Name          string    `json:"name"`
	Env           string    `json:"env"`
	SourcePath    string    `json:"source_path"`
}

type UpdateConfigurationRequest struct {
	ApplicationID    uuid.UUID  `json:"application_id"`
	Name             string     `json:"name"`
	Env              string     `json:"env"`
	SourcePath       string     `json:"source_path"`
	LatestRevisionNo int        `json:"latest_revision_no"`
	LatestRevisionID *uuid.UUID `json:"latest_revision_id,omitempty"`
}

func NewConfigurationHandler() *ConfigurationHandler {
	return &ConfigurationHandler{
		svc: app.ConfigurationService,
	}
}

// Create
// @Summary 创建配置
// @Description 创建一个新的配置
// @Tags Configuration
// @Accept json
// @Produce json
// @Param data body api.CreateConfigurationRequest true "Configuration Data"
// @Success 201 {object} httpx.DataResponse[domain.Configuration]
// @Router /api/v1/configurations [post]
func (h *ConfigurationHandler) Create(c *gin.Context) {
	var req CreateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	cfg := &domain.Configuration{
		ApplicationID: req.ApplicationID,
		Name:          req.Name,
		Env:           req.Env,
		SourcePath:    req.SourcePath,
	}
	cfg.WithCreateDefault()

	_, err := h.svc.Create(c.Request.Context(), cfg)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteData(c, http.StatusCreated, cfg)
}

// Get
// @Summary 获取配置
// @Tags    Configuration
// @Param   id path string true "Configuration ID"
// @Success 200 {object} httpx.DataResponse[domain.Configuration]
// @Router  /api/v1/configurations/{id} [get]
func (h *ConfigurationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	cfg, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
		return
	}

	httpx.WriteData(c, http.StatusOK, cfg)
}

// Update
// @Summary 更新配置
// @Tags    Configuration
// @Param   id   path string               true "Configuration ID"
// @Param   data body api.UpdateConfigurationRequest true "Configuration Data"
// @Success 204
// @Router  /api/v1/configurations/{id} [put]
func (h *ConfigurationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	var req UpdateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	cfg := domain.Configuration{
		ApplicationID:    req.ApplicationID,
		Name:             req.Name,
		Env:              req.Env,
		SourcePath:       req.SourcePath,
		LatestRevisionNo: req.LatestRevisionNo,
		LatestRevisionID: req.LatestRevisionID,
	}
	cfg.SetID(id)

	if err := h.svc.Update(c.Request.Context(), &cfg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteNoContent(c)
}

// Delete
// @Summary 删除配置
// @Tags    Configuration
// @Param   id path string true "Configuration ID"
// @Success 204
// @Router  /api/v1/configurations/{id} [delete]
func (h *ConfigurationHandler) Delete(c *gin.Context) {
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

// List
// @Summary 获取配置列表
// @Tags    Configuration
// @Success 200 {object} httpx.ListResponse[domain.Configuration]
// @Router  /api/v1/configurations [get]
func (h *ConfigurationHandler) List(c *gin.Context) {
	filter := app.ConfigurationListFilter{
		IncludeDeleted: httpx.IncludeDeleted(c),
		Name:           c.Query("name"),
	}

	cfgs, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	total := len(cfgs)
	cfgs = httpx.PaginateSlice(cfgs, paging)
	httpx.WriteList(c, http.StatusOK, cfgs, paging, total)
}

// Sync
// @Summary 从集中配置仓同步配置 revision
// @Tags    Configuration
// @Param   id path string true "Configuration ID"
// @Success 200 {object} httpx.DataResponse[domain.ConfigurationRevision]
// @Router  /api/v1/configurations/{id}/sync [post]
func (h *ConfigurationHandler) Sync(c *gin.Context) {
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
		case errors.Is(err, app.ErrConfigSourceNotFound):
			httpx.WriteError(c, http.StatusFailedDependency, "failed_precondition", err.Error(), nil)
		case errors.Is(err, app.ErrConfigRepositoryUnavailable):
			httpx.WriteError(c, http.StatusFailedDependency, "failed_precondition", err.Error(), nil)
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		}
		return
	}

	httpx.WriteData(c, http.StatusOK, result.Revision)
}
