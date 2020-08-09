package herodote

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
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
ORDER BY
  date DESC
LIMIT $1
OFFSET $2
`
const searchCommitQueryWithVector = `
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
  search_vector @@ to_tsquery('english', $3)
ORDER BY
  date DESC
LIMIT $1
OFFSET $2
`

func (a app) searchCommit(ctx context.Context, query string, page, pageSize uint) ([]Commit, uint, error) {
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

	if len(words) == 0 {
		return list, totalCount, db.List(ctx, a.db, scanner, searchCommitQuery, pageSize, (page-1)*pageSize)
	}

	return list, totalCount, db.List(ctx, a.db, scanner, searchCommitQueryWithVector, pageSize, (page-1)*pageSize, strings.Join(words, " | "))
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
