DROP MATERIALIZED VIEW IF EXISTS herodote.lexeme;
DROP INDEX IF EXISTS words;

CREATE INDEX commit_search ON herodote.commit USING gist(search_vector);
