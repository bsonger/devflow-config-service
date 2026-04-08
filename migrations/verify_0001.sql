DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'configurations') THEN
    RAISE EXCEPTION 'missing table: configurations';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'configuration_revisions') THEN
    RAISE EXCEPTION 'missing table: configuration_revisions';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'configurations' AND column_name = 'source_path'
  ) THEN
    RAISE EXCEPTION 'missing column: configurations.source_path';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'configuration_revisions' AND column_name = 'source_commit'
  ) THEN
    RAISE EXCEPTION 'missing column: configuration_revisions.source_commit';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'configuration_revisions' AND column_name = 'source_digest'
  ) THEN
    RAISE EXCEPTION 'missing column: configuration_revisions.source_digest';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = 'public' AND indexname = 'uq_configuration_revisions_no'
  ) THEN
    RAISE EXCEPTION 'missing index: uq_configuration_revisions_no';
  END IF;
END $$;
