--- clean
DROP MATERIALIZED VIEW IF EXISTS herodote.lexeme;

DROP TABLE IF EXISTS herodote.commit;
DROP TABLE IF EXISTS herodote.repository;

DROP INDEX IF EXISTS words;
DROP INDEX IF EXISTS repository_id;
DROP INDEX IF EXISTS commit_component;
DROP INDEX IF EXISTS commit_type;

DROP SCHEMA IF EXISTS herodote;

-- extension
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- schema
CREATE SCHEMA herodote;

-- repository
CREATE SEQUENCE herodote.repository_seq;
CREATE TABLE herodote.repository (
  id BIGINT NOT NULL DEFAULT nextval('herodote.repository_seq'),
  remote TEXT NOT NULL,
  name TEXT NOT NULL
);
ALTER SEQUENCE herodote.repository_seq OWNED BY herodote.repository.id;

CREATE UNIQUE INDEX repository_id ON herodote.repository(id);

-- commit
CREATE TABLE herodote.commit (
  hash TEXT NOT NULL,
  repository_id BIGINT NOT NULL REFERENCES herodote.repository(id) ON DELETE CASCADE,
  component TEXT NOT NULL,
  type TEXT NOT NULL,
  content TEXT NOT NULL,
  search_vector TSVECTOR,
  date TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX commit_component ON herodote.commit(component);
CREATE INDEX commit_type ON herodote.commit(type);

-- lexeme
CREATE MATERIALIZED VIEW herodote.lexeme AS
  SELECT word FROM ts_stat('SELECT search_vector FROM herodote.commit');

CREATE INDEX words ON herodote.lexeme USING gin(word gin_trgm_ops);

-- data
DO $$
  DECLARE repository_id BIGINT;

  BEGIN
    INSERT INTO herodote.repository (remote, name) VALUES ('github.com', 'ViBiOh/herodote') RETURNING id INTO repository_id;
    INSERT INTO herodote.commit (hash, repository_id, component, type, content, date) VALUES
      ('f2aab4b', repository_id, 'style', '', 'Prettifying files', to_timestamp(1596655426)),
      ('950bdd1', repository_id, 'chore', '', 'Setting appropriate timeout for curl command', to_timestamp(1596655426)),
      ('df7cb61', repository_id, 'chore', 'deps', 'bump @testing-library/user-event from 12.0.15 to 12.1.0', to_timestamp(1596655426)),
      ('8514506', repository_id, 'chore', 'deps', 'bump @testing-library/jest-dom from 5.11.1 to 5.11.2', to_timestamp(1596655426)),
      ('6234821', repository_id, 'chore', 'deps-dev', 'bump eslint-plugin-react from 7.20.4 to 7.20.5', to_timestamp(1596655426)),
      ('3be6f46', repository_id, 'chore', 'deps', 'bump @testing-library/user-event from 7.2.1 to 12.0.15', to_timestamp(1596655426)),
      ('bdfc820', repository_id, 'chore', 'deps-dev', 'bump eslint-plugin-react from 7.20.3 to 7.20.4', to_timestamp(1596655426)),
      ('ff5e2f1', repository_id, 'style', '', 'Formating yaml, markdown and json files', to_timestamp(1596655426)),
      ('5b6cec5', repository_id, 'style', '', 'Formating yaml, markdown and json files', to_timestamp(1596655426)),
      ('30368a7', repository_id, 'feat', '', 'Adding throbber on load-more button', to_timestamp(1596655426)),
      ('7bffc4a', repository_id, 'chore', 'deps', 'bump funtch from 2.1.0 to 2.1.1', to_timestamp(1596655426)),
      ('c3c6f45', repository_id, 'chore', 'deps', 'bump @testing-library/react from 9.5.0 to 10.4.7', to_timestamp(1596655426)),
      ('6ba4616', repository_id, 'chore', 'deps', 'bump algoliasearch from 4.3.0 to 4.3.1', to_timestamp(1596655426)),
      ('dab7d96', repository_id, 'feat', '', 'Adding powered by algolia logo', to_timestamp(1596655426)),
      ('9017728', repository_id, 'feat', '', 'Adding favicon and OpenGraph', to_timestamp(1596655426)),
      ('df85112', repository_id, 'docs', '', 'Improving docs with conventional commits checker and behavior', to_timestamp(1596655426)),
      ('7d6d5d8', repository_id, 'docs', '', 'Improving README for getting started', to_timestamp(1596655426)),
      ('958469d', repository_id, 'refactor', '', 'Moving URL read and query to App instead of Filters', to_timestamp(1596655426)),
      ('207b214', repository_id, 'fix', 'a11y', 'Adding aria-label for input text', to_timestamp(1596655426)),
      ('e6efa1b', repository_id, 'feat', '', 'Adding date separator on rendering', to_timestamp(1596655426)),
      ('a5091e9', repository_id, 'ci', '', 'Changing k8s limit to have QoS Guaranteed', to_timestamp(1596655426)),
      ('1600d05', repository_id, 'feat', '', 'Sorting facets by alpha value', to_timestamp(1596655426)),
      ('707cfd5', repository_id, 'fix', '', 'Adding z-index for filters__values on mobile', to_timestamp(1596655426)),
      ('c065e22', repository_id, 'feat', '', 'Handling responsive design', to_timestamp(1596655426)),
      ('6ec5545', repository_id, 'refactor', '', 'Splitting containers in two separate entity', to_timestamp(1596655426)),
      ('207a605', repository_id, 'feat', '', 'Adding throbber and error management', to_timestamp(1596655426))
    ;

    UPDATE herodote.commit SET search_vector = to_tsvector('english', unaccent(hash)) || to_tsvector('english', unaccent(component)) || to_tsvector('english', unaccent(type)) || to_tsvector('english', unaccent(content));

    REFRESH MATERIALIZED VIEW herodote.lexeme;
END $$;

SELECT coalesce((string_agg(word, ' | ')), '') FROM herodote.lexeme WHERE similarity(word, unaccent('react')) > 0.2;
SELECT * FROM herodote.commit WHERE search_vector @@ to_tsquery('english', 'refactor | readm | read | react | eslint-plugin-react | /react');
