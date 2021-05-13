package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	httpModel "github.com/ViBiOh/httputils/v4/pkg/model"
	"github.com/lib/pq"
)

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

func (a app) SearchCommit(ctx context.Context, query string, filters map[string][]string, before, after string, page, pageSize uint) ([]model.Commit, uint, error) {
	var words []string
	var err error

	if len(query) > 0 {
		words, err = a.findSimilarWords(ctx, query)
		if err != nil {
			return nil, 0, httpModel.WrapNotFound(fmt.Errorf("unable to find similar words: %s", err))
		}

		if len(words) == 0 {
			return nil, 0, httpModel.WrapNotFound(errors.New("query doesn't match any known words"))
		}
	}

	var totalCount uint
	var list []model.Commit

	scanner := func(rows *sql.Rows) error {
		var item model.Commit

		if err := rows.Scan(&item.Hash, &item.Type, &item.Component, &item.Revert, &item.Breaking, &item.Content, &item.Date, &item.Remote, &item.Repository, &totalCount); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	}

	sqlQuery, sqlArgs := computeSearchQuery(page, pageSize, words, filters, before, after)

	return list, totalCount, db.List(ctx, a.db, scanner, sqlQuery, sqlArgs...)
}

func computeSearchQuery(page, pageSize uint, words []string, filters map[string][]string, before, after string) (string, []interface{}) {
	query := strings.Builder{}
	query.WriteString(searchCommitQuery)

	args := []interface{}{
		pageSize,
		(page - 1) * pageSize,
	}

	if len(words) != 0 {
		args = append(args, strings.Join(words, " | "))
		query.WriteString(fmt.Sprintf(" AND search_vector @@ to_tsquery('english', $%d)", len(args)))
	}

	for key, values := range filters {
		if len(values) == 0 {
			continue
		}

		sqlValues := make([]string, 0)
		for _, value := range values {
			if len(value) == 0 {
				continue
			}

			sqlValues = append(sqlValues, strings.ToLower(value))
		}

		if len(sqlValues) == 0 {
			continue
		}

		args = append(args, pq.Array(sqlValues))
		query.WriteString(fmt.Sprintf(" AND %s = ANY($%d)", key, len(args)))
	}

	if len(before) != 0 {
		args = append(args, before)
		query.WriteString(fmt.Sprintf(" AND date < $%d", len(args)))
	}

	if len(after) != 0 {
		args = append(args, after)
		query.WriteString(fmt.Sprintf(" AND date > $%d", len(args)))
	}

	query.WriteString(searchCommitTail)

	return query.String(), args
}

const findSimilarWordsQuery = `
SELECT DISTINCT
  word
FROM
  herodote.lexeme
WHERE
  similarity(word, unaccent($1)) > 0.4
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
