package api

import (
	"net/http"

	"github.com/bsonger/devflow-config-service/pkg/model"
	"github.com/bsonger/devflow-config-service/pkg/service"
	"github.com/bsonger/devflow-service-common/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ConfigurationRouteApi = NewConfigurationHandler()

type ConfigurationHandler struct {
}

func NewConfigurationHandler() *ConfigurationHandler {
	return &ConfigurationHandler{}
}

// Create
// @Summary 创建配置
// @Description 创建一个新的配置
// @Tags Configuration
// @Accept json
// @Produce json
// @Param data body model.Configuration true "Configuration Data"
// @Success 200 {object} httpx.CreateResponse
// @Router /api/v1/configurations [post]
func (h *ConfigurationHandler) Create(c *gin.Context) {
	var cfg *model.Configuration
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg.WithCreateDefault()

	id, err := service.ConfigurationService.Create(c.Request.Context(), cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpx.CreateResponse{ID: id.String()})
}

// Get
// @Summary 获取配置
// @Tags    Configuration
// @Param   id path string true "Configuration ID"
// @Success 200 {object} model.Configuration
// @Router  /api/v1/configurations/{id} [get]
func (h *ConfigurationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cfg, err := service.ConfigurationService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

// Update
// @Summary 更新配置
// @Tags    Configuration
// @Param   id   path string               true "Configuration ID"
// @Param   data body model.Configuration true "Configuration Data"
// @Success 200  {object} map[string]string
// @Router  /api/v1/configurations/{id} [put]
func (h *ConfigurationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var cfg model.Configuration
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg.SetID(id)

	if err := service.ConfigurationService.Update(c.Request.Context(), &cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete
// @Summary 删除配置
// @Tags    Configuration
// @Param   id path string true "Configuration ID"
// @Success 200 {object} map[string]string
// @Router  /api/v1/configurations/{id} [delete]
func (h *ConfigurationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := service.ConfigurationService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// List
// @Summary 获取配置列表
// @Tags    Configuration
// @Success 200 {array} model.Configuration
// @Router  /api/v1/configurations [get]
func (h *ConfigurationHandler) List(c *gin.Context) {
	filter := primitive.M{}
	if !httpx.IncludeDeleted(c) {
		filter["deleted_at"] = primitive.M{"$exists": false}
	}
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}

	cfgs, err := service.ConfigurationService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := len(cfgs)
	cfgs = httpx.PaginateSlice(cfgs, paging)
	httpx.SetPaginationHeaders(c, total, paging)

	c.JSON(http.StatusOK, cfgs)
}
