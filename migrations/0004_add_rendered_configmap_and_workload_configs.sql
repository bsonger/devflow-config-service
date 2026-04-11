ALTER TABLE configuration_revisions
  ADD COLUMN IF NOT EXISTS rendered_configmap JSONB NOT NULL DEFAULT '{"data":{}}'::jsonb;

CREATE TABLE IF NOT EXISTS workload_configs (
  id UUID PRIMARY KEY,
  application_id UUID NOT NULL,
  environment_id TEXT NULL,
  name TEXT NOT NULL,
  replicas INTEGER NOT NULL DEFAULT 1,
  resources JSONB NOT NULL DEFAULT '{}'::jsonb,
  probes JSONB NOT NULL DEFAULT '{}'::jsonb,
  env JSONB NOT NULL DEFAULT '[]'::jsonb,
  workload_type TEXT NOT NULL DEFAULT '',
  strategy TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_workload_configs_scope_name_active
  ON workload_configs (application_id, COALESCE(environment_id, ''), name)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_workload_configs_application_id
  ON workload_configs (application_id);
