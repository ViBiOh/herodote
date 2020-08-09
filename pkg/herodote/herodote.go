package herodote

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/cron"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/httpjson"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/query"
	"github.com/ViBiOh/httputils/v3/pkg/request"
)

const (
	commitsPath = "/commits"
	filtersPath = "/filters"
	refreshPath = "/refresh"
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

func (a app) Start() {
	cron.New().Days().At("06:00").In("Europe/Paris").Start(a.refreshLexeme, func(err error) {
		logger.Error("unable to refresh lexeme: %s", err)
	})
}

// Handler for request. Should be use with net/http
func (a app) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Authorization") != a.secret {
			httperror.Unauthorized(w, ErrAuthentificationFailed)
			return
		}

		if strings.HasPrefix(r.URL.Path, commitsPath) {
			a.handleCommits(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, filtersPath) {
			a.handleFilters(w, r)
			return
		}

		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, refreshPath) {
			if err := a.refreshLexeme(time.Now()); err != nil {
				httperror.InternalServerError(w, err)
				return
			}

			return
		}

		httperror.NotFound(w)
	})
}

func (a app) handleCommits(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.handlePostCommits(w, r)
	} else if r.Method == http.MethodGet {
		a.handleGetCommits(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a app) handlePostCommits(w http.ResponseWriter, r *http.Request) {
	data, err := request.ReadBodyRequest(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	commit = commit.Sanitize()
	if err := commit.Check(); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if err := a.saveCommit(r.Context(), commit); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a app) handleGetCommits(w http.ResponseWriter, r *http.Request) {
	pagination, err := query.ParsePagination(r, 1, 20, 100)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	params := r.URL.Query()

	query := strings.TrimSpace(params.Get("q"))
	filters := map[string][]string{
		"repository": params["repository"],
		"type":       params["type"],
		"component":  params["component"],
	}

	commits, totalCount, err := a.searchCommit(r.Context(), query, filters, pagination.Page, pagination.PageSize)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	httpjson.ResponsePaginatedJSON(w, http.StatusOK, pagination.Page, pagination.PageSize, totalCount, commits, httpjson.IsPretty(r))
}

func (a app) handleFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	name := r.URL.Query().Get("name")

	var values []string
	var err error

	switch strings.ToLower(name) {
	case "repository":
		values, err = a.listRepositories(ctx)
	case "type":
		values, err = a.listTypes(ctx)
	case "component":
		values, err = a.listComponents(ctx)
	default:
		httperror.NotFound(w)
		return
	}

	if err != nil {
		httperror.InternalServerError(w, err)
	}

	httpjson.ResponseArrayJSON(w, http.StatusOK, values, httpjson.IsPretty(r))
}