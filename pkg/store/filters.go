package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const listRepositoryFiltersQuery = `
SELECT DISTINCT
  repository
FROM
  herodote.commit
ORDER BY
  repository
`

const listTypeFiltersQuery = `
SELECT DISTINCT
  type
FROM
  herodote.commit
ORDER BY
  type
`

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

func (a app) ListFilters(ctx context.Context, name string) ([]string, error) {
	var list []string

	scanner := func(rows *sql.Rows) error {
		var item string
		if err := rows.Scan(&item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	}

	var query string

	switch name {
	case "repository":
		query = listRepositoryFiltersQuery
	case "type":
		query = listTypeFiltersQuery
	case "component":
		query = listComponentFiltersQuery
	default:
		return nil, fmt.Errorf("unknown filter `%s`", name)
	}

	return list, db.List(ctx, a.db, scanner, query)
}
