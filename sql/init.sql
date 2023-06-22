CREATE TABLE IF NOT EXISTS todolist (
  id bytea NOT NULL,
  title text NOT NULL,
  created_at timestamptz NOT NULL,
  done_at timestamptz,
  is_done BOOLEAN GENERATED ALWAYS AS (done_at IS NOT NULL) STORED,

  PRIMARY KEY(id)
);
