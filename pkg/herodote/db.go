package herodote

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/lib/pq"
)

const insertCommitQuery = `
INSERT INTO
  herodote.commit
(
  hash,
  type,
  component,
  revert,
  breaking,
  content,
  date,
  remote,
  repository,
  search_vector
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  to_timestamp($7),
  $8,
  $9,
  to_tsvector('english', $1) || to_tsvector('english', $2) || to_tsvector('english', $3) || to_tsvector('english', $6)
)
`

func (a app) saveCommit(ctx context.Context, o Commit) error {
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		return db.Exec(ctx, insertCommitQuery, o.Hash, o.Type, o.Component, o.Revert, o.Breaking, o.Content, o.Date, o.Remote, o.Repository)
	})
}

const searchCommitQuery = `
SELECT
  hash,
  type,
  component,
  revert,
  breaking,
  content,
  date,
  remote,
  repository,
  count(1) OVER() AS full_count
FROM
  herodote.commit
WHERE
  TRUE
`

const searchCommitTail = `
ORDER BY
  date DESC
LIMIT $1
OFFSET $2
`

func (a app) searchCommit(ctx context.Context, query string, filters map[string][]string, page, pageSize uint) ([]Commit, uint, error) {
	words, err := a.findSimilarWords(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to find similar words: %s", err)
	}

	var totalCount uint
	var list []Commit

	scanner := func(rows *sql.Rows) error {
		var item Commit
		var date time.Time

		if err := rows.Scan(&item.Hash, &item.Type, &item.Component, &item.Revert, &item.Breaking, &item.Content, &date, &item.Remote, &item.Repository, &totalCount); err != nil {
			return err
		}

		item.Date = uint64(date.Unix())
		list = append(list, item)
		return nil
	}

	sqlQuery, sqlArgs := computeSearchQuery(page, pageSize, words, filters)

	return list, totalCount, db.List(ctx, a.db, scanner, sqlQuery, sqlArgs...)
}

func computeSearchQuery(page, pageSize uint, words []string, filters map[string][]string) (string, []interface{}) {
	query := searchCommitQuery
	args := []interface{}{
		pageSize,
		(page - 1) * pageSize,
	}

	if len(words) != 0 {
		args = append(args, strings.Join(words, " | "))
		query += fmt.Sprintf(" AND search_vector @@ to_tsquery('english', $%d)", len(args))
	}

	for key, values := range filters {
		if len(values) == 0 {
			continue
		}

		sqlValues := make([]string, len(values))
		for index, value := range values {
			sqlValues[index] = strings.ToLower(value)
		}

		args = append(args, pq.Array(sqlValues))
		query += fmt.Sprintf(" AND %s = ANY($%d)", key, len(args))
	}

	query += searchCommitTail

	return query, args
}

const findSimilarWordsQuery = `
SELECT DISTINCT
  word
FROM
  herodote.lexeme
WHERE
  similarity(word, unaccent($1)) > 0.2
`

func (a app) findSimilarWords(ctx context.Context, query string) ([]string, error) {
	var list []string

	scanner := func(rows *sql.Rows) error {
		var item string
		if err := rows.Scan(&item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	}

	return list, db.List(ctx, a.db, scanner, findSimilarWordsQuery, query)
}

func (a app) listFilters(ctx context.Context, query string) ([]string, error) {
	var list []string

	scanner := func(rows *sql.Rows) error {
		var item string
		if err := rows.Scan(&item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	}

	return list, db.List(ctx, a.db, scanner, query)
}

const listRepositoryFiltersQuery = `
SELECT DISTINCT
  repository
FROM
  herodote.commit
ORDER BY
  repository
`

func (a app) listRepositories(ctx context.Context) ([]string, error) {
	return a.listFilters(ctx, listRepositoryFiltersQuery)
}

const listTypeFiltersQuery = `
SELECT DISTINCT
  type
FROM
  herodote.commit
ORDER BY
  type
`

func (a app) listTypes(ctx context.Context) ([]string, error) {
	return a.listFilters(ctx, listTypeFiltersQuery)
}

const listComponentFiltersQuery = `
SELECT DISTINCT
  component
FROM
  herodote.commit
WHERE
  component <> ''
ORDER BY
  component
`

func (a app) listComponents(ctx context.Context) ([]string, error) {
	return a.listFilters(ctx, listComponentFiltersQuery)
}