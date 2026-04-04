# Architecture

## Purpose

`devflow-config-service` 负责配置元数据面，只管理 `Configuration`。

## Inbound Surface

- `GET/POST/PUT/DELETE /api/v1/configurations`

## Data And Dependencies

- 主存储：MongoDB
- 主要集合：`configurations`
- 启动、路由、HTTP 公共件、观测基础设施来自 `devflow-service-common`

## Outbound Rules

- 当前主流程没有必须的跨服务 RPC
- 如果后续增加调用其他服务或外部系统，必须同时产生 `metrics + trace + structured log`

## Non-Goals

- 不负责 `Project`
- 不负责 `Application`
- 不负责 `Manifest`、`Job`、`Intent`
- 不负责 verify webhook / writeback
