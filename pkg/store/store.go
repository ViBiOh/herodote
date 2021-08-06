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
type App struct {
	db db.App
}

// New creates new App from Config
func New(db db.App) App {
	return App{
		db: db,
	}
}

// Enabled checks if storage is ok
func (a App) Enabled() bool {
	return a.db.Enabled()
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

// SaveCommit into the storage
func (a App) SaveCommit(ctx context.Context, o model.Commit) error {
	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		return a.db.Exec(ctx, insertCommitQuery, o.Hash, o.Type, o.Component, o.Revert, o.Breaking, o.Content, o.Date.Unix(), o.Remote, o.Repository)
	})
}

const refreshFiltersQuery = `REFRESH MATERIALIZED VIEW herodote.filters`

// Refresh filters values
func (a App) Refresh(ctx context.Context) error {
	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		return a.db.Exec(ctx, refreshFiltersQuery)
	})
}
