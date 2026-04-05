# Configuration

## Ownership

- owner repo: `devflow-config-service`
- authoritative model file: `pkg/model/configuration.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/swagger.yaml` (transitional; still reflects legacy handler layer until API migration)

## Purpose

`Configuration` 是发布配置元数据资源，供 release 路径消费。
实际可变内容已经拆分到不可变的 `ConfigurationRevision`。

## Common base fields

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `id` | `uuid.UUID` | server-generated | no | 主键 |
| `created_at` | `time.Time` | server-generated | no | 创建时间 |
| `updated_at` | `time.Time` | server-generated | no | 更新时间 |
| `deleted_at` | `*time.Time` | optional | system-managed | 软删除时间 |

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `application_id` | `uuid.UUID` | required | user | 所属应用 ID |
| `name` | `string` | required | user | 配置名 |
| `env` | `string` | required | user | 目标环境 |
| `latest_revision_no` | `int` | system-managed | no | 当前最新 revision 序号 |
| `latest_revision_id` | `*uuid.UUID` | optional/system-managed | no | 当前最新 revision ID |

## Nested types

## Related child resource: `ConfigurationRevision`

| Field | Type | Description |
|---|---|---|
| `configuration_id` | `uuid.UUID` | 所属配置 |
| `revision_no` | `int` | 版本号 |
| `files` | `[]File` | 配置文件快照 |
| `content_hash` | `string` | 内容哈希 |
| `message` | `string` | 变更说明 |
| `created_by` | `string` | 创建人 |
| `created_at` | `time.Time` | 创建时间 |

### `File`
- `name: string`
- `content: string`

## Create / update rules

### Create
- target relational contract:
  - required: `application_id`, `name`, `env`
  - create flow should also materialize the first immutable revision
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - `name`, `env`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
  - `latest_revision_no`, `latest_revision_id`

## Validation notes

- `application_id` 必须引用存在的 Application
- revision 一旦创建不可修改
- 列表默认过滤掉软删除数据

## Source pointers

- router: `pkg/router/configuration.go`
- handler: `pkg/api/configuration.go`
- service: `pkg/service/configuration.go`
- model: `pkg/model/configuration.go`
