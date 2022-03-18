package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/jackc/pgx/v4"
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

// SearchCommit in the storage based on given filters
func (a App) SearchCommit(ctx context.Context, query string, filters map[string][]string, before, after string, pageSize uint, last string) ([]model.Commit, uint, error) {
	var words []string
	if len(query) > 0 {
		words = strings.Split(query, " ")
	}

	var totalCount uint
	var list []model.Commit

	scanner := func(rows pgx.Rows) error {
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

func computeSearchQuery(pageSize uint, last string, words []string, filters map[string][]string, before, after string) (string, []any) {
	query := strings.Builder{}
	query.WriteString(searchCommitQuery)

	args := []any{
		pageSize,
	}

	if len(words) != 0 {
		args = append(args, strings.Join(words, " & "))
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

		args = append(args, sqlValues)
		query.WriteString(fmt.Sprintf(" AND %s = ANY($%d)", key, len(args)))
	}

	args = computeDateQuery(&query, args, before, last, after)

	query.WriteString(searchCommitTail)

	return query.String(), args
}

func computeDateQuery(query *strings.Builder, args []any, before, last, after string) []any {
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
