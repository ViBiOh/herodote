package main

import (
	"flag"
	"os"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/herodote/pkg/herodote"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/cors"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/owasp"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/redis"
	"github.com/ViBiOh/httputils/v4/pkg/renderer"
	"github.com/ViBiOh/httputils/v4/pkg/server"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
)

type configuration struct {
	appServer  server.Config
	promServer server.Config
	health     health.Config
	alcotest   alcotest.Config
	logger     logger.Config
	tracer     tracer.Config
	prometheus prometheus.Config
	owasp      owasp.Config
	cors       cors.Config
	renderer   renderer.Config
	herodote   herodote.Config
	db         db.Config
	redis      redis.Config
}

func newConfig() (configuration, error) {
	fs := flag.NewFlagSet("herodote", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	return configuration{
		appServer:  server.Flags(fs, ""),
		promServer: server.Flags(fs, "prometheus", flags.NewOverride("Port", uint(9090)), flags.NewOverride("IdleTimeout", 10*time.Second), flags.NewOverride("ShutdownTimeout", 5*time.Second)),
		health:     health.Flags(fs, ""),
		alcotest:   alcotest.Flags(fs, ""),
		logger:     logger.Flags(fs, "logger"),
		tracer:     tracer.Flags(fs, "tracer"),
		prometheus: prometheus.Flags(fs, "prometheus", flags.NewOverride("Gzip", false)),
		owasp:      owasp.Flags(fs, "", flags.NewOverride("Csp", "default-src 'self'; base-uri 'self'; script-src 'self' 'httputils-nonce'; style-src 'self' 'httputils-nonce'")),
		cors:       cors.Flags(fs, "cors"),
		renderer:   renderer.Flags(fs, "", flags.NewOverride("Title", "Herodote"), flags.NewOverride("PublicURL", "https://herodote.vibioh.fr")),
		herodote:   herodote.Flags(fs, ""),
		db:         db.Flags(fs, "db"),
		redis:      redis.Flags(fs, "redis"),
	}, fs.Parse(os.Args[1:])
}
