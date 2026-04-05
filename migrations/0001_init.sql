CREATE TABLE IF NOT EXISTS configurations (
  id UUID PRIMARY KEY,
  application_id UUID NOT NULL,
  name TEXT NOT NULL,
  env TEXT NOT NULL,
  latest_revision_no INTEGER NOT NULL DEFAULT 1,
  latest_revision_id UUID NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_configurations_app_env_name_active
  ON configurations (application_id, env, name)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_configurations_application_id
  ON configurations (application_id);

CREATE TABLE IF NOT EXISTS configuration_revisions (
  id UUID PRIMARY KEY,
  configuration_id UUID NOT NULL,
  revision_no INTEGER NOT NULL,
  files JSONB NOT NULL DEFAULT '[]'::jsonb,
  content_hash TEXT NOT NULL,
  message TEXT NOT NULL DEFAULT '',
  created_by TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_configuration_revisions_no
  ON configuration_revisions (configuration_id, revision_no);

CREATE INDEX IF NOT EXISTS idx_configuration_revisions_hash
  ON configuration_revisions (content_hash);
