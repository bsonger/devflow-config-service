package router

import (
	"github.com/bsonger/devflow-config-service/pkg/api"
	"github.com/gin-gonic/gin"
)

func RegisterWorkloadConfigRoutes(rg *gin.RouterGroup) {
	workloads := rg.Group("/workload-configs")
	workloads.GET("", api.WorkloadConfigRouteApi.List)
	workloads.GET("/:id", api.WorkloadConfigRouteApi.Get)
	workloads.POST("", api.WorkloadConfigRouteApi.Create)
	workloads.PUT("/:id", api.WorkloadConfigRouteApi.Update)
	workloads.DELETE("/:id", api.WorkloadConfigRouteApi.Delete)
}
