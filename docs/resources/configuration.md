# Configuration

## Ownership

- owner repo: `devflow-config-service`
- authoritative model file: `pkg/model/configuration.go`
- authoritative API doc: `docs/api-spec.md`
- swagger source: `docs/swagger.yaml`

## Purpose

`Configuration` 是发布配置元数据资源，供 release 路径消费。

## Common base fields

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `id` | `ObjectID` | server-generated | no | 主键 |
| `created_at` | `time.Time` | server-generated | no | 创建时间 |
| `updated_at` | `time.Time` | server-generated | no | 更新时间 |
| `deleted_at` | `*time.Time` | optional | system-managed | 软删除时间 |

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `name` | `string` | expected on create | user | 配置名 |
| `files` | `[]*File` | optional | user | 配置文件集合 |

## Nested types

### `File`
- `name: string`
- `content: string`

## Create / update rules

### Create
- current API behavior:
  - handler 绑定整个 `model.Configuration`
  - 当前未做额外字段级 `binding:"required"` 校验
- practical required fields:
  - `name`
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - `name`, `files`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`

## Validation notes

- 路径参数 `id` 必须是合法 ObjectID
- 列表默认过滤掉软删除数据
- 当前实现没有更细的字段级校验规则文档化

## Source pointers

- router: `pkg/router/configuration.go`
- handler: `pkg/api/configuration.go`
- service: `pkg/service/configuration.go`
- model: `pkg/model/configuration.go`
