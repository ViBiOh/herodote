package main

import (
	"embed"
	"flag"
	"os"

	"github.com/ViBiOh/herodote/pkg/herodote"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/cors"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/httputils"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/owasp"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/recoverer"
	"github.com/ViBiOh/httputils/v4/pkg/renderer"
	"github.com/ViBiOh/httputils/v4/pkg/server"
)

//go:embed templates static
var content embed.FS

func main() {
	fs := flag.NewFlagSet("herodote", flag.ExitOnError)

	appServerConfig := server.Flags(fs, "")
	promServerConfig := server.Flags(fs, "prometheus", flags.NewOverride("Port", 9090), flags.NewOverride("IdleTimeout", "10s"), flags.NewOverride("ShutdownTimeout", "5s"))
	healthConfig := health.Flags(fs, "")

	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus", flags.NewOverride("Gzip", false))
	owaspConfig := owasp.Flags(fs, "", flags.NewOverride("Csp", "default-src 'self'; base-uri 'self'; script-src 'self' 'nonce'; style-src 'self' 'nonce'"))
	corsConfig := cors.Flags(fs, "cors")
	rendererConfig := renderer.Flags(fs, "", flags.NewOverride("Title", "Herodote"), flags.NewOverride("PublicURL", "https://herodote.vibioh.fr"))

	herodoteConfig := herodote.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	appServer := server.New(appServerConfig)
	promServer := server.New(promServerConfig)
	prometheusApp := prometheus.New(prometheusConfig)

	herodoteDb, err := db.New(dbConfig)
	logger.Fatal(err)
	defer herodoteDb.Close()

	healthApp := health.New(healthConfig, herodoteDb.Ping)

	storeApp := store.New(herodoteDb)
	herodoteApp, err := herodote.New(herodoteConfig, storeApp)
	logger.Fatal(err)

	rendererApp, err := renderer.New(rendererConfig, content, herodote.FuncMap)
	logger.Fatal(err)

	rendererHandler := rendererApp.Handler(herodoteApp.TemplateFunc)

	go promServer.Start("prometheus", healthApp.End(), prometheusApp.Handler())
	go appServer.Start("http", healthApp.End(), httputils.Handler(rendererHandler, healthApp, recoverer.Middleware, prometheusApp.Middleware, owasp.New(owaspConfig).Middleware, cors.New(corsConfig).Middleware))

	healthApp.WaitForTermination(appServer.Done())
	server.GracefulWait(appServer.Done(), promServer.Done())
}
