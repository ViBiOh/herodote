package main

import (
	"context"
	"fmt"

	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/redis"
	"github.com/ViBiOh/httputils/v4/pkg/request"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
)

type clients struct {
	logger     *logger.Logger
	prometheus *prometheus.App
	tracer     tracer.App
	redis      redis.Client
	database   db.App
	health     *health.App
}

func newClients(ctx context.Context, config configuration) (clients, error) {
	var output clients
	var err error

	output.logger = logger.New(config.logger)
	logger.Global(output.logger)

	output.prometheus = prometheus.New(config.prometheus)

	output.tracer, err = tracer.New(ctx, config.tracer)
	if err != nil {
		return output, fmt.Errorf("tracer: %w", err)
	}

	request.AddTracerToDefaultClient(output.tracer.GetProvider())

	output.database, err = db.New(ctx, config.db, output.tracer.GetTracer("database"))
	if err != nil {
		return output, fmt.Errorf("database: %w", err)
	}

	output.redis, err = redis.New(config.redis, output.tracer.GetProvider())
	if err != nil {
		return output, fmt.Errorf("redis: %w", err)
	}

	output.health = health.New(config.health, output.database.Ping)

	return output, nil
}

func (c clients) Close(ctx context.Context) {
	c.database.Close()
	c.redis.Close()
	c.tracer.Close(ctx)
	c.logger.Close()
}
