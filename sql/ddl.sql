--- clean
DROP MATERIALIZED VIEW IF EXISTS herodote.filters;

DROP TABLE IF EXISTS herodote.commit;

DROP INDEX IF EXISTS words;
DROP INDEX IF EXISTS commit_repository;
DROP INDEX IF EXISTS commit_component;
DROP INDEX IF EXISTS commit_type;

DROP SCHEMA IF EXISTS herodote;

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
CREATE INDEX commit_search ON herodote.commit USING gist(search_vector);

-- filters
CREATE MATERIALIZED VIEW herodote.filters (
  kind,
  value
) AS
  SELECT DISTINCT 'repository', repository FROM herodote.commit
  UNION SELECT DISTINCT 'type', type FROM herodote.commit
  UNION SELECT DISTINCT 'component', component FROM herodote.commit WHERE component IS NOT NULL;
