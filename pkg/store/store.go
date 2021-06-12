package store

import (
	"context"
	"errors"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/httputils/v4/pkg/db"
)

var (
	// ErrNotFound occurs when nothing found
	ErrNotFound = errors.New("not found")
)

// App of package
type App interface {
	SaveCommit(context.Context, model.Commit) error
	SearchCommit(context.Context, string, map[string][]string, string, string, uint, string) ([]model.Commit, uint, error)
	ListFilters(context.Context) (map[string][]string, error)
	Refresh(context.Context) error
}

type app struct {
	db db.App
}

// New creates new App from Config
func New(db db.App) App {
	return app{
		db: db,
	}
}

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

func (a app) SaveCommit(ctx context.Context, o model.Commit) error {
	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		return a.db.Exec(ctx, insertCommitQuery, o.Hash, o.Type, o.Component, o.Revert, o.Breaking, o.Content, o.Date.Unix(), o.Remote, o.Repository)
	})
}

const refreshLexemeQuery = `REFRESH MATERIALIZED VIEW herodote.lexeme`
const refreshFiltersQuery = `REFRESH MATERIALIZED VIEW herodote.filters`

func (a app) Refresh(ctx context.Context) error {
	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		err := a.db.Exec(ctx, refreshLexemeQuery)
		if err != nil {
			return err
		}

		err = a.db.Exec(ctx, refreshFiltersQuery)
		return err
	})
}
