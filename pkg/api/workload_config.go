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

var WorkloadConfigRouteApi = NewWorkloadConfigHandler()

type workloadConfigService interface {
	Create(ctx context.Context, item *domain.WorkloadConfig) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.WorkloadConfig, error)
	Update(ctx context.Context, item *domain.WorkloadConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter app.WorkloadConfigListFilter) ([]domain.WorkloadConfig, error)
}

type WorkloadConfigHandler struct{ svc workloadConfigService }

func NewWorkloadConfigHandler() *WorkloadConfigHandler {
	return &WorkloadConfigHandler{svc: app.WorkloadConfigService}
}

// @Summary Create workload config
// @Tags WorkloadConfig
// @Accept json
// @Produce json
// @Param data body domain.WorkloadConfigInput true "WorkloadConfig data"
// @Success 201 {object} httpx.DataResponse[domain.WorkloadConfig]
// @Router /api/v1/workload-configs [post]
func (h *WorkloadConfigHandler) Create(c *gin.Context) {
	var req domain.WorkloadConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	item := &domain.WorkloadConfig{
		ApplicationID: req.ApplicationID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Replicas:      req.Replicas,
		Resources:     req.Resources,
		Probes:        req.Probes,
		Env:           req.Env,
		WorkloadType:  req.WorkloadType,
		Strategy:      req.Strategy,
	}
	item.WithCreateDefault()
	if _, err := h.svc.Create(c.Request.Context(), item); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	httpx.WriteData(c, http.StatusCreated, item)
}

// @Summary Get workload config
// @Tags WorkloadConfig
// @Produce json
// @Param id path string true "WorkloadConfig ID"
// @Success 200 {object} httpx.DataResponse[domain.WorkloadConfig]
// @Router /api/v1/workload-configs/{id} [get]
func (h *WorkloadConfigHandler) Get(c *gin.Context) {
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

// @Summary Update workload config
// @Tags WorkloadConfig
// @Accept json
// @Param id path string true "WorkloadConfig ID"
// @Param data body domain.WorkloadConfigInput true "WorkloadConfig data"
// @Success 204
// @Router /api/v1/workload-configs/{id} [put]
func (h *WorkloadConfigHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}
	var req domain.WorkloadConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}
	item := &domain.WorkloadConfig{
		ApplicationID: req.ApplicationID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Replicas:      req.Replicas,
		Resources:     req.Resources,
		Probes:        req.Probes,
		Env:           req.Env,
		WorkloadType:  req.WorkloadType,
		Strategy:      req.Strategy,
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

// @Summary Delete workload config
// @Tags WorkloadConfig
// @Param id path string true "WorkloadConfig ID"
// @Success 204
// @Router /api/v1/workload-configs/{id} [delete]
func (h *WorkloadConfigHandler) Delete(c *gin.Context) {
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

// @Summary List workload configs
// @Tags WorkloadConfig
// @Produce json
// @Param application_id query string false "Application ID"
// @Param environment_id query string false "Environment ID"
// @Param name query string false "Name"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpx.ListResponse[domain.WorkloadConfig]
// @Router /api/v1/workload-configs [get]
func (h *WorkloadConfigHandler) List(c *gin.Context) {
	var filter app.WorkloadConfigListFilter
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
