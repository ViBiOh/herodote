package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/db"
)

var (
	// ErrNotFound occurs when nothing found
	ErrNotFound = errors.New("not found")
)

// App of package
type App interface {
	SaveCommit(context.Context, model.Commit) error
	SearchCommit(context.Context, string, map[string][]string, string, string, uint, uint) ([]model.Commit, uint, error)
	ListFilters(context.Context, string) ([]string, error)
	Refresh(context.Context) error
}

type app struct {
	db *sql.DB
}

// New creates new App from Config
func New(db *sql.DB) App {
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
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		return db.Exec(ctx, insertCommitQuery, o.Hash, o.Type, o.Component, o.Revert, o.Breaking, o.Content, o.Date, o.Remote, o.Repository)
	})
}

const refreshMaterializedView = `REFRESH MATERIALIZED VIEW herodote.lexeme`

func (a app) Refresh(ctx context.Context) error {
	_, err := a.db.ExecContext(context.Background(), refreshMaterializedView)
	return err
}
