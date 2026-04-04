package main

import (
	"github.com/bsonger/devflow-config-service/pkg/router"
	"github.com/bsonger/devflow-config-service/platform/shared/bootstrap"
)

func main() {
	err := bootstrap.Run(bootstrap.Options{
		Name: "config-service",
		RouteOptions: router.Options{
			ServiceName:   "config-service",
			EnableSwagger: true,
			Modules: []router.Module{
				router.ModuleConfiguration,
			},
		},
		PortEnv:        "CONFIG_SERVICE_PORT",
		DefaultPort:    8082,
		MetricsPortEnv: "CONFIG_SERVICE_METRICS_PORT",
		PprofPortEnv:   "CONFIG_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
