CREATE TABLE records (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  data_key     TEXT       NOT NULL,
  url          TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  type_key     TEXT       NOT NULL,
  content      TEXT       NOT NULL,
  history      TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_storage_key   ON records(data_key);

CREATE INDEX idx_storage_title ON records(type_key, title);

