package herodote

import (
	"context"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const insertCommitQuery = `
INSERT INTO
  herodote.commit
(
  hash,
  type,
  component,
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
  to_timestamp($5),
  $6,
  $7,
  to_tsvector('english', $1) || to_tsvector('english', $2) || to_tsvector('english', $3) || to_tsvector('english', $4)
)
`

func (a app) saveCommit(ctx context.Context, o Commit) error {
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		return db.Exec(ctx, insertCommitQuery, o.Hash, o.Type, o.Component, o.Content, o.Date, o.Remote, o.Repository)
	})
}
