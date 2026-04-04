# API Spec

## Resources

- `Configuration`
  - `GET /api/v1/configurations`
  - `POST /api/v1/configurations`
  - `GET /api/v1/configurations/:id`
  - `PUT /api/v1/configurations/:id`
  - `DELETE /api/v1/configurations/:id`

## Response Rules

- 创建接口返回统一创建响应
- 列表接口遵循统一分页参数和分页响应头

## Error Rules

- 非法 ObjectID 返回 `400`
- 资源不存在返回 `404`
- 存储层或未分类错误返回 `500`

## Swagger

Swagger 必须只包含 `Configuration` 接口。
