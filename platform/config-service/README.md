# Config Service

职责：

- 提供 `Configuration` 的 CRUD
- 作为配置元数据的单独服务边界

当前实现：

- `cmd/main.go` 通过 `devflow-service-common/bootstrap` 启动
- `pkg/api/configuration.go`
- `pkg/service/configuration.go`
- `pkg/router/configuration.go`
- `pkg/config/config.go`

建议端口：

- `CONFIG_SERVICE_PORT`
- `CONFIG_SERVICE_METRICS_PORT`
- `CONFIG_SERVICE_PPROF_PORT`

运行时：

- 上报的 OTel `service.name` 为 `config-service`
- 任何 outbound service / external call 都必须带 `metrics + trace + structured log`
- 默认 harness 为 `Planner -> Generator -> Evaluator`，并且支持 delegation 时必须真实启动 sub-agents
