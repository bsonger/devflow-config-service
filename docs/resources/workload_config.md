# Workload Config

## Ownership

- owner repo: `devflow-config-service`
- authoritative model file: `pkg/domain/workload_config.go`
- authoritative API doc: `docs/api-spec.md`

## Purpose

`WorkloadConfig` 是运行模板资源。
它描述副本数、资源限制、探针、环境变量、工作负载类型，以及当前阶段允许的策略类型。

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `application_id` | `uuid.UUID` | required | user | 所属应用 ID |
| `environment_id` | `string` | optional | user | 可选环境级覆盖 |
| `name` | `string` | required | user | 配置名 |
| `replicas` | `int` | required | user | 副本数 |
| `resources` | `object` | optional | user | 资源配额 |
| `probes` | `object` | optional | user | 探针配置 |
| `env` | `object[]` | optional | user | 环境变量 |
| `workload_type` | `string` | required | user | 工作负载类型 |
| `strategy` | `string` | optional | user | 发布策略类型 |

## Rules

- `WorkloadConfig` 只提供 CRUD
- 不提供独立 validate 接口
- 策略当前只保留类型，不包含步骤明细

## Source pointers

- router: `pkg/router/workload_config.go`
- handler: `pkg/api/workload_config.go`
- app: `pkg/app/workload_config.go`
- model: `pkg/domain/workload_config.go`
