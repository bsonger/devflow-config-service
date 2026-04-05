# Configuration

## Ownership

- owner repo: `devflow-config-service`
- authoritative model file: `pkg/domain/configuration.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/generated/swagger/swagger.yaml`

## Purpose

`Configuration` 是发布配置元数据资源，供 release 路径消费。
实际可变内容存放在集中配置仓里，这个资源只维护逻辑身份、`source_path` 和最新 revision 指针。

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
| `source_path` | `string` | required | user | 集中配置仓中的稳定路径 |
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
| `source_commit` | `string` | 来源仓的冻结 commit/ref |
| `source_digest` | `string` | 源文件集合摘要 |
| `message` | `string` | 变更说明 |
| `created_by` | `string` | 创建人 |
| `created_at` | `time.Time` | 创建时间 |

### `File`
- `name: string`
- `content: string`

## Create / update rules

### Create
- target relational contract:
  - required: `application_id`, `name`, `env`, `source_path`
  - create flow only creates identity, not revision content
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - `name`, `env`, `source_path`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
  - `latest_revision_no`, `latest_revision_id`

### Sync
- `POST /api/v1/configurations/{id}/sync`
- 从集中配置仓 `source_path/files` 读取当前内容
- 内容没变化时返回当前最新 revision
- 内容变化时创建新的不可变 `ConfigurationRevision`

## Validation notes

- `application_id` 必须引用存在的 Application
- `source_path` 必须映射到集中配置仓中的有效目录
- revision 一旦创建不可修改
- 列表默认过滤掉软删除数据

## Source pointers

- router: `pkg/router/configuration.go`
- handler: `pkg/api/configuration.go`
- app: `pkg/app/configuration.go`
- sync flow: `pkg/app/configuration_sync.go`
- repo reader: `pkg/infra/config_repo/repository.go`
- model: `pkg/domain/configuration.go`
