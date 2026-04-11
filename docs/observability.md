# Observability

## Mandatory Rule

只要涉及调用其他服务或外部系统，就必须同时产出：

- metrics
- trace
- structured log

## Inbound HTTP

- 所有业务请求必须有 server span
- 记录请求次数、耗时、错误数
- `/metrics`、`/healthz`、`/readyz`、`/debug/pprof/*` 不计入业务指标

## Outbound

- 当前 config-service 主路径没有必须的跨服务调用
- 若后续增加出站调用，必须补齐 client span、调用计数、延迟、错误计数和结构化日志

## Log Fields

- 基础字段：`service`、`trace_id`、`span_id`、`request_id`
- 资源字段：`app_config_id`、`workload_config_id`

## Profile

- `pprof` 和 Pyroscope 都是显式开启的诊断能力
- 当前仓没有额外 repo-local telemetry 扩展
