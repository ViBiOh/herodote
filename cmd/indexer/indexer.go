package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("indexer", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	ctx := context.Background()

	herodoteDb, err := db.New(ctx, dbConfig, nil)
	logger.Fatal(err)
	defer herodoteDb.Close()

	logger.Info("Lexeme refresh...")
	logger.Fatal(store.New(herodoteDb).Refresh(ctx))
	logger.Info("Lexeme refreshed!")
}
