package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/herodote/pkg/herodote"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/cors"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
)

func main() {
	fs := flag.NewFlagSet("herodote", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

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

	go herodoteApp.Start()
	httputils.New(serverConfig).ListenAndServe(herodoteApp.Handler(), []model.Middleware{
		prometheus.New(prometheusConfig).Middleware,
		owasp.New(owaspConfig).Middleware,
		cors.New(corsConfig).Middleware,
	}, herodoteDb.Ping)
}
