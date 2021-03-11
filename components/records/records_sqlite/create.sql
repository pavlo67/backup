DROP TABLE IF EXISTS records;

CREATE TABLE records (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  urn          TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  type_key     TEXT       NOT NULL,
  data         TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,

  owner_id     TEXT       NOT NULL,
  viewer_id    TEXT       NOT NULL,
  history      TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_records_urn          ON records(urn);
CREATE INDEX idx_records_viewer_title ON records(viewer_id, type_key, title);
CREATE INDEX idx_records_owner_title  ON records(owner_id,  type_key, title);
CREATE INDEX idx_records_type_title   ON records(           type_key, title);

