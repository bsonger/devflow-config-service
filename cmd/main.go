package main

import (
	_ "github.com/bsonger/devflow-config-service/docs/generated/swagger"
	"github.com/bsonger/devflow-config-service/pkg/infra/config"
	"github.com/bsonger/devflow-config-service/pkg/router"
	"github.com/bsonger/devflow-service-common/bootstrap"
	"github.com/bsonger/devflow-service-common/observability"
)

func main() {
	err := bootstrap.Run(bootstrap.Options[config.Config, router.Options, string]{
		Name:         "config-service",
		RouteOptions: router.Options{ServiceName: "config-service", EnableSwagger: true},
		Load:         config.Load,
		InitRuntime:  config.InitRuntime,
		NewRouter: func(opts router.Options) bootstrap.Runner {
			return router.NewRouterWithOptions(opts)
		},
		ResolveConfigPort:  config.ResolveConfigPort,
		StartMetricsServer: observability.StartMetricsServer,
		StartPprofServer:   observability.StartPprofServer,
		PortEnv:            "CONFIG_SERVICE_PORT",
		DefaultPort:        8082,
		MetricsPortEnv:     "CONFIG_SERVICE_METRICS_PORT",
		PprofPortEnv:       "CONFIG_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
