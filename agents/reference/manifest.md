# Manifest Reference Redirect

This file is index-only.

Authoritative sources:

- `devflow-control/docs/resources/manifest.md`
- `devflow-control/docs/services/release-service.md`
- `devflow-release-service/docs/resources/manifest.md`

This file must not define manifest semantics.
- 状态枚举：
  - `Pending`
  - `Running`
  - `Succeeded`
  - `Failed`
- `steps` 记录各任务步骤的执行状态与时间戳
- 当前确认的创建链路：
  - `pkg/api/manifest.go` 通过 `POST /api/v1/manifests` 调用 `ManifestService.CreateManifest`
  - `intent` 模式下，先落库 Manifest，再创建 build intent，并返回 `execution_intent_id`
  - `direct` 模式下，直接提交 Tekton PipelineRun，再把 `pipeline_id`、`steps` 等元数据写回 Manifest

## Must Not

- 不把 `Manifest` 当作 Tekton 真相来源的替代品
- 不把本地创建成功写成构建成功

## Outputs

- Manifest 元数据
- build 过程的 `steps`
- 与 build intent 的关联

## Pass/Fail

- `Pass`：Manifest 语义仍是构建快照和状态读模型
- `Fail`：Manifest 变成执行器私有状态或终态来源被写乱
