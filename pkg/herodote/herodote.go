package herodote

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v3/pkg/cron"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/httpjson"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/query"
	rendererModel "github.com/ViBiOh/httputils/v3/pkg/renderer/model"
	"github.com/ViBiOh/httputils/v3/pkg/request"
)

const (
	isoDateLayout = "2006-01-02"
	dayDuration   = time.Hour * 24

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
	TemplateFunc(*http.Request) (string, int, map[string]interface{}, error)
	Start()
}

// Config of package
type Config struct {
	secret *string
}

type app struct {
	store  store.App
	secret string
	colors map[string]string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		secret: flags.New(prefix, "herodote").Name("HttpSecret").Default("").Label("HTTP Secret Key for Update").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, store store.App) (App, error) {
	secret := strings.TrimSpace(*config.secret)
	if len(secret) == 0 {
		return nil, errors.New("http secret is required")
	}

	if store == nil {
		return nil, errors.New("store is required")
	}

	return app{
		secret: secret,
		store:  store,
		colors: make(map[string]string),
	}, nil
}

func (a app) Start() {
	cron.New().Days().At("06:00").In("Europe/Paris").Start(func(_ time.Time) error {
		logger.Info("Refreshing lexeme")
		return a.store.Refresh(context.Background())
	}, func(err error) {
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

		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, refreshPath) {
			if err := a.store.Refresh(r.Context()); err != nil {
				httperror.InternalServerError(w, err)
				return
			}

			return
		}

		httperror.NotFound(w)
	})
}

func (a app) TemplateFunc(r *http.Request) (string, int, map[string]interface{}, error) {
	commits, _, _, err := a.listCommits(r)
	if err != nil {
		return "", http.StatusInternalServerError, nil, err
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to parse query: %s", err)
	}

	repositories, err := a.store.ListFilters(r.Context(), "repository")
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to list repositories: %s", err)
	}

	types, err := a.store.ListFilters(r.Context(), "type")
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to list types: %s", err)
	}

	components, err := a.store.ListFilters(r.Context(), "component")
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to list components: %s", err)
	}

	return "public", http.StatusOK, map[string]interface{}{
		"Path":         r.URL.Path,
		"Filters":      params,
		"Repositories": repositories,
		"Types":        types,
		"Components":   components,
		"Commits":      commits,
	}, nil
}

func (a app) listCommits(r *http.Request) ([]model.Commit, uint, query.Pagination, error) {
	pagination, err := query.ParsePagination(r, 1, 50, 100)
	if err != nil {
		return nil, 0, pagination, rendererModel.WrapInvalid(err)
	}

	params := r.URL.Query()

	query := strings.TrimSpace(params.Get("q"))
	filters := map[string][]string{
		"repository": params["repository"],
		"type":       params["type"],
		"component":  params["component"],
	}

	before := strings.TrimSpace(params.Get("before"))
	if err := checkDate(before); err != nil {
		return nil, 0, pagination, rendererModel.WrapInvalid(err)
	}

	after := strings.TrimSpace(params.Get("after"))
	if err := checkDate(after); err != nil {
		return nil, 0, pagination, rendererModel.WrapInvalid(err)
	}

	commits, totalCount, err := a.store.SearchCommit(r.Context(), query, filters, before, after, pagination.Page, pagination.PageSize)
	return commits, totalCount, pagination, err
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

func (a app) handleGetCommits(w http.ResponseWriter, r *http.Request) {
	commits, totalCount, pagination, err := a.listCommits(r)
	if err != nil {
		if errors.Is(err, rendererModel.ErrInvalid) {
			httperror.BadRequest(w, err)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	httpjson.ResponsePaginatedJSON(w, http.StatusOK, pagination.Page, pagination.PageSize, totalCount, commits, httpjson.IsPretty(r))
}

func (a app) handlePostCommits(w http.ResponseWriter, r *http.Request) {
	data, err := request.ReadBodyRequest(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	var commit model.Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	commit = commit.Sanitize()
	if err := commit.Check(); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if err := a.store.SaveCommit(r.Context(), commit); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func checkDate(raw string) error {
	if len(raw) == 0 {
		return nil
	}

	_, err := time.Parse(isoDateLayout, raw)
	if err != nil {
		return fmt.Errorf("unable to parse date: %s", err)
	}

	return nil
}
