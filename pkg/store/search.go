package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ViBiOh/herodote/pkg/model"
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
`

func (a app) SearchCommit(ctx context.Context, query string, filters map[string][]string, before, after string, pageSize uint, last string) ([]model.Commit, uint, error) {
	var words []string
	var err error

	if len(query) > 0 {
		words, err = a.findSimilarWords(ctx, query)
		if err != nil {
			return nil, 0, httpModel.WrapNotFound(fmt.Errorf("unable to find similar words: %s", err))
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

	sqlQuery, sqlArgs := computeSearchQuery(pageSize, last, words, filters, before, after)

	return list, totalCount, a.db.List(ctx, scanner, sqlQuery, sqlArgs...)
}

func computeSearchQuery(pageSize uint, last string, words []string, filters map[string][]string, before, after string) (string, []interface{}) {
	query := strings.Builder{}
	query.WriteString(searchCommitQuery)

	args := []interface{}{
		pageSize,
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

	args = computeDateQuery(&query, args, before, last, after)

	query.WriteString(searchCommitTail)

	return query.String(), args
}

func computeDateQuery(query *strings.Builder, args []interface{}, before, last, after string) []interface{} {
	if len(before) != 0 || len(last) != 0 {
		if len(last) != 0 {
			args = append(args, last)
		} else {
			args = append(args, before)
		}

		query.WriteString(fmt.Sprintf(" AND date < $%d", len(args)))
	}

	if len(after) != 0 {
		args = append(args, after)
		query.WriteString(fmt.Sprintf(" AND date > $%d", len(args)))
	}

	return args
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

	return list, a.db.List(ctx, scanner, findSimilarWordsQuery, query)
}
