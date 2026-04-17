ALTER TABLE configurations
  ADD COLUMN IF NOT EXISTS mount_path TEXT NOT NULL DEFAULT '/etc/devflow/config';
