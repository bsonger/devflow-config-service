# Constraints

## Ownership

- `AppConfig` 和 `WorkloadConfig` 属于 config-service 独占边界

## Prohibited

- 不得在本仓库新增 `Project`、`Application`、`Manifest`、`Release`、`Intent`、`Verify` 对外 API
- 不得把 release / verify 运行态逻辑塞回 config-service
- 不得在 metrics label 中写入高基数业务主键

## Data Rules

- 删除语义沿用现有模型定义
- 任何写操作都必须更新 `updated_at`
