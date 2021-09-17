package store

import (
	"context"

	"github.com/jackc/pgx/v4"
)

const listFiltersQuery = `
SELECT
  kind,
  value
FROM
  herodote.filters
`

// ListFilters available on GUI
func (a App) ListFilters(ctx context.Context) (map[string][]string, error) {
	list := make(map[string][]string)

	scanner := func(rows pgx.Rows) error {
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
