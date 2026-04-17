# App Config

## Ownership

- owner repo: `devflow-config-service`
- authoritative model file: `pkg/domain/app_config.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/generated/swagger/swagger.yaml`

## Purpose

`AppConfig` 是发布配置元数据资源，供 release 路径消费。
实际文件内容来自固定配置仓，这个资源只维护逻辑身份、派生后的 `source_path` 和最新 revision 指针。当前 `source_path` 规则为 `applications/devflow-platform/services/<name>`。
它还维护 `mount_path`，用于声明后续渲染 Deployment 时应把配置文件挂到容器内哪个目录或文件路径。

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
| `environment_id` | `string` | required | user | 目标环境 |
| `mount_path` | `string` | optional | user | 配置文件在容器内的挂载目标，默认 `/etc/devflow/config` |
| `source_path` | `string` | system-derived | system | 固定配置仓中的派生路径 |
| `latest_revision_no` | `int` | system-managed | no | 当前最新 revision 序号 |
| `latest_revision_id` | `*uuid.UUID` | optional/system-managed | no | 当前最新 revision ID |

## Nested types

## Related child resource: `AppConfigRevision`

| Field | Type | Description |
|---|---|---|
| `app_config_id` | `uuid.UUID` | 所属配置 |
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
  - required: `application_id`, `environment_id`, `name`
  - optional: `mount_path`
  - create flow only creates identity, not revision content
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - `name`, `environment_id`, `mount_path`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
  - `latest_revision_no`, `latest_revision_id`

### Sync
- `POST /api/v1/app-configs/{id}/sync-from-repo`
- 固定从 `git@github.com:bsonger/devflow-config-repo.git` 的 `main` 分支读取
- 在冻结 revision 前先执行一次 `origin/main` 的快进拉取
- 路径由 `name` 推导为 `applications/devflow-platform/services/<name>`
- 默认冻结 `configuration.yaml`
- 当 `environment_id` 不是空 / `base` 且存在 `environments/<environment_id>.yaml` 时，一并冻结该 overlay 文件
- `deployment.yaml` / `service.yaml` 属于 workload / network 侧输入，不属于 `AppConfig` 冻结内容
- 内容没变化时返回当前最新 revision
- 内容变化时创建新的不可变 `AppConfigRevision`

## Validation notes

- `application_id` 必须引用存在的 Application
- `source_path` 必须映射到固定配置仓中的有效目录
- revision 一旦创建不可修改
- 列表默认过滤掉软删除数据

## Source pointers

- router: `pkg/router/app_config.go`
- handler: `pkg/api/app_config.go`
- app: `pkg/app/app_config.go`
- repo reader: `pkg/infra/config_repo/repository.go`
- model: `pkg/domain/app_config.go`
