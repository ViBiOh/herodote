package herodote

import (
	"database/sql"
	"errors"
	"flag"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/v3/pkg/cron"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

var (
	// ErrAuthentificationFailed occurs when secret is invalid
	ErrAuthentificationFailed = errors.New("invalid secret provided")
)

// App of package
type App interface {
	Handler() http.Handler
	Start()
}

// Config of package
type Config struct {
	secret *string
}

type app struct {
	secret string

	db *sql.DB
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		secret: flags.New(prefix, "herodote").Name("HttpSecret").Default("").Label("HTTP Secret Key for Update").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, database *sql.DB) (App, error) {
	secret := strings.TrimSpace(*config.secret)
	if len(secret) == 0 {
		return nil, errors.New("http secret is required")
	}

	if database == nil {
		return nil, errors.New("database is required")
	}

	return app{
		secret: secret,
		db:     database,
	}, nil

}

// Handler for request. Should be use with net/http
func (a app) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Authorization") != a.secret {
			httperror.Unauthorized(w, ErrAuthentificationFailed)
		}
	})
}

func (a app) Start() {
	cron.New().Days().At("06:00").In("Europe/Paris").Start(a.refreshLexeme, func(err error) {
		logger.Error("unable to refresh lexeme: %s", err)
	})
}
