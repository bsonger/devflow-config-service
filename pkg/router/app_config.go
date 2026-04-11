package router

import (
	"github.com/bsonger/devflow-config-service/pkg/api"
	"github.com/gin-gonic/gin"
)

func RegisterAppConfigRoutes(rg *gin.RouterGroup) {
	appConfigs := rg.Group("/app-configs")
	appConfigs.GET("", api.AppConfigRouteApi.List)
	appConfigs.GET("/:id", api.AppConfigRouteApi.Get)
	appConfigs.POST("", api.AppConfigRouteApi.Create)
	appConfigs.PUT("/:id", api.AppConfigRouteApi.Update)
	appConfigs.DELETE("/:id", api.AppConfigRouteApi.Delete)
	appConfigs.POST("/:id/sync-from-repo", api.AppConfigRouteApi.Sync)
}
