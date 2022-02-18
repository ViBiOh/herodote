package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
)

func main() {
	fs := flag.NewFlagSet("indexer", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	herodoteDb, err := db.New(dbConfig, tracer.App{})
	logger.Fatal(err)
	defer herodoteDb.Close()

	logger.Info("Lexeme refresh...")
	logger.Fatal(store.New(herodoteDb).Refresh(context.Background()))
	logger.Info("Lexeme refreshed!")
}
