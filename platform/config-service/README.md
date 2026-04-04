# Config Service

职责：

- 提供 `Configuration` 的 CRUD
- 作为配置元数据的单独服务边界

当前复用的现有实现：

- `pkg/api/configuration.go`
- `pkg/service/configuration.go`
- `pkg/router/configuration.go`

建议端口：

- `CONFIG_SERVICE_PORT`
- `CONFIG_SERVICE_METRICS_PORT`
- `CONFIG_SERVICE_PPROF_PORT`

运行时：

- 上报的 OTel `service.name` 为 `config-service`
