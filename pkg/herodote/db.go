package herodote

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const insertCommitQuery = `
INSERT INTO
  herodote.commit
(
  hash,
  repository_id,
  type,
  component,
  content,
  date,
  search_vector
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  to_timestamp($6),
  to_tsvector('english', $1) || to_tsvector('english', $3) || to_tsvector('english', $4) || to_tsvector('english', $5)
)
`

func (a app) saveCommit(ctx context.Context, o Commit) error {
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		repositoryID, err := a.getOrCreateRepository(ctx, o.Remote, o.Repository)
		if err != nil {
			return fmt.Errorf("unable to get or create repository for %s: %s", o.Repository, err)
		}

		return db.Exec(ctx, insertCommitQuery, o.Hash, repositoryID, o.Type, o.Component, o.Content, o.Date)
	})
}

const getRepositoryIDQuery = `
SELECT
  id
FROM
  herodote.repository
WHERE
  name = $1
`

const insertRepositoryQuery = `
INSERT INTO
  herodote.repository
(
  remote,
  name
) VALUES (
  $1,
  $2
) RETURNING ID
`

func (a app) getOrCreateRepository(ctx context.Context, remote, name string) (uint64, error) {
	var id uint64
	scanner := func(row *sql.Row) error {
		return row.Scan(&id)
	}

	err := db.Get(ctx, a.db, scanner, getRepositoryIDQuery, name)
	if err == nil {
		return id, nil
	}

	if err == sql.ErrNoRows {
		return db.Create(ctx, insertRepositoryQuery, remote, name)
	}

	return id, err
}
