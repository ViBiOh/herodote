--- clean
DROP MATERIALIZED VIEW IF EXISTS herodote.filters;
DROP MATERIALIZED VIEW IF EXISTS herodote.lexeme;

DROP TABLE IF EXISTS herodote.commit;

DROP INDEX IF EXISTS words;
DROP INDEX IF EXISTS commit_repository;
DROP INDEX IF EXISTS commit_component;
DROP INDEX IF EXISTS commit_type;

DROP SCHEMA IF EXISTS herodote;

-- extension
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- schema
CREATE SCHEMA herodote;

-- commit
CREATE TABLE herodote.commit (
  repository TEXT NOT NULL,
  hash TEXT NOT NULL,
  type TEXT NOT NULL,
  component TEXT NOT NULL,
  revert BOOLEAN NOT NULL,
  breaking BOOLEAN NOT NULL,
  content TEXT NOT NULL,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  remote TEXT NOT NULL,
  search_vector TSVECTOR
);

CREATE UNIQUE INDEX commit_id ON herodote.commit(repository, hash);
CREATE INDEX commit_repository ON herodote.commit(repository);
CREATE INDEX commit_component ON herodote.commit(component);
CREATE INDEX commit_type ON herodote.commit(type);

-- lexeme
CREATE MATERIALIZED VIEW herodote.lexeme AS
  SELECT word FROM ts_stat('SELECT search_vector FROM herodote.commit');

CREATE INDEX words ON herodote.lexeme USING gin(word gin_trgm_ops);

-- filters
CREATE MATERIALIZED VIEW herodote.filters (
  kind,
  value
) AS
  SELECT DISTINCT 'repository', repository FROM herodote.commit
  UNION SELECT DISTINCT 'type', type FROM herodote.commit
  UNION SELECT DISTINCT 'component', component FROM herodote.commit WHERE component <> '';
