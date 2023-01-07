package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/herodote/pkg/version"
	"github.com/ViBiOh/httputils/v4/pkg/cache"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/redis"
	"github.com/ViBiOh/httputils/v4/pkg/sha"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
)

type App struct {
	redis redis.App
	store store.App
}

func New(redis redis.App, database db.App) App {
	return App{
		redis: redis,
		store: store.New(database),
	}
}

func (a App) Enabled() bool {
	return a.store.Enabled()
}

func (a App) ListFilters(ctx context.Context) (map[string][]string, error) {
	return a.store.ListFilters(ctx)
}

func (a App) SearchCommit(ctx context.Context, query string, filters map[string][]string, before, after string, pageSize uint, last string) (model.CommitsList, error) {
	searchHash := sha.Stream().Write(query).Write(filters).Write(before).Write(after).Write(pageSize).Write(last).Sum()

	return cache.Load(ctx, a.redis, version.Redis("commits:"+searchHash), func(ctx context.Context) (model.CommitsList, error) {
		return a.store.SearchCommit(ctx, query, filters, before, after, pageSize, last)
	}, time.Hour)
}

func (a App) SaveCommit(ctx context.Context, commit model.Commit) error {
	err := a.store.SaveCommit(ctx, commit)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	go func(ctx context.Context) {
		if err := a.redis.DeletePattern(ctx, version.Redis("commits:*")); err != nil {
			logger.Error("redis delete after save commit: %s", err)
		}
	}(tracer.CopyToBackground(ctx))

	return err
}
