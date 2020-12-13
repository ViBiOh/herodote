package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/ViBiOh/herodote/pkg/herodote"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/cors"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
	"github.com/ViBiOh/httputils/v3/pkg/renderer"
)

const (
	apiPath = "/api"
)

func main() {
	fs := flag.NewFlagSet("herodote", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "", flags.NewOverride("Csp", "default-src 'self'; base-uri 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"))
	corsConfig := cors.Flags(fs, "cors")
	rendererConfig := renderer.Flags(fs, "", flags.NewOverride("Title", "Herodote"), flags.NewOverride("PublicURL", "https://herodote.vibioh.fr"))

	herodoteConfig := herodote.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	herodoteDb, err := db.New(dbConfig)
	logger.Fatal(err)

	storeApp := store.New(herodoteDb)
	herodoteApp, err := herodote.New(herodoteConfig, storeApp)
	logger.Fatal(err)

	rendererApp, err := renderer.New(rendererConfig, herodote.FuncMap)
	logger.Fatal(err)

	herodoteHandler := http.StripPrefix(apiPath, herodoteApp.Handler())
	rendererHandler := rendererApp.Handler(herodoteApp.TemplateFunc)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, apiPath) {
			herodoteHandler.ServeHTTP(w, r)
			return
		}

		rendererHandler.ServeHTTP(w, r)
	})

	go herodoteApp.Start()
	httputils.New(serverConfig).ListenAndServe(handler, []model.Pinger{herodoteDb.Ping}, prometheus.New(prometheusConfig).Middleware, owasp.New(owaspConfig).Middleware, cors.New(corsConfig).Middleware)
}
