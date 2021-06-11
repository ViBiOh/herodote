package store

import (
	"context"
	"database/sql"
)

const listFiltersQuery = `
SELECT
  kind,
  value
FROM
  herodote.filters
`

func (a app) ListFilters(ctx context.Context) (map[string][]string, error) {
	list := make(map[string][]string)

	scanner := func(rows *sql.Rows) error {
		var kind, value string
		if err := rows.Scan(&kind, &value); err != nil {
			return err
		}

		if values, ok := list[kind]; !ok {
			list[kind] = []string{value}
		} else {
			list[kind] = append(values, value)
		}

		return nil
	}

	return list, a.db.List(ctx, scanner, listFiltersQuery)
}
