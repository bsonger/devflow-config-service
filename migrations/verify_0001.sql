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
    WHERE table_name = 'configuration_revisions' AND column_name = 'env_vars'
  ) THEN
    RAISE EXCEPTION 'missing column: configuration_revisions.env_vars';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = 'public' AND indexname = 'uq_configuration_revisions_no'
  ) THEN
    RAISE EXCEPTION 'missing index: uq_configuration_revisions_no';
  END IF;
END $$;
