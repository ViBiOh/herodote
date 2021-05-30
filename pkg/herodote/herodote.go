package herodote

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/httperror"
	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	httpModel "github.com/ViBiOh/httputils/v4/pkg/model"
	"github.com/ViBiOh/httputils/v4/pkg/query"
)

const (
	isoDateLayout = "2006-01-02"

	apiPath     = "/api"
	commitsPath = "/commits"
)

var (
	// ErrAuthentificationFailed occurs when secret is invalid
	ErrAuthentificationFailed = errors.New("invalid secret provided")
)

// App of package
type App interface {
	Handler() http.Handler
	TemplateFunc(http.ResponseWriter, *http.Request) (string, int, map[string]interface{}, error)
}

// Config of package
type Config struct {
	secret *string
}

type app struct {
	apiHandler http.Handler
	colors     map[string]string
	storeApp   store.App
	secret     string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		secret: flags.New(prefix, "herodote").Name("HttpSecret").Default("").Label("HTTP Secret Key for Update").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, storeApp store.App) (App, error) {
	secret := strings.TrimSpace(*config.secret)
	if len(secret) == 0 {
		return nil, errors.New("http secret is required")
	}

	if storeApp == nil {
		return nil, errors.New("store is required")
	}

	app := app{
		secret:   secret,
		storeApp: storeApp,
		colors:   make(map[string]string),
	}

	app.apiHandler = http.StripPrefix(apiPath, app.Handler())

	return app, nil
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

		httperror.NotFound(w)
	})
}

func (a app) TemplateFunc(w http.ResponseWriter, r *http.Request) (string, int, map[string]interface{}, error) {
	if strings.HasPrefix(r.URL.Path, apiPath) {
		a.apiHandler.ServeHTTP(w, r)
		return "", 0, nil, nil
	}

	commits, _, _, err := a.listCommits(r)
	if err != nil {
		return "", http.StatusInternalServerError, nil, err
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to parse query: %s", err)
	}

	filters, err := a.storeApp.ListFilters(r.Context())
	if err != nil {
		return "", http.StatusInternalServerError, nil, fmt.Errorf("unable to list filters: %s", err)
	}

	return "public", http.StatusOK, map[string]interface{}{
		"Path":         r.URL.Path,
		"Filters":      params,
		"Repositories": filters["repository"],
		"Types":        filters["type"],
		"Components":   filters["component"],
		"Commits":      commits,
		"Now":          time.Now(),
	}, nil
}

func (a app) listCommits(r *http.Request) ([]model.Commit, uint, query.Pagination, error) {
	pagination, err := query.ParsePagination(r, 50, 100)
	if err != nil {
		return nil, 0, pagination, httpModel.WrapInvalid(err)
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
		return nil, 0, pagination, httpModel.WrapInvalid(err)
	}

	after := strings.TrimSpace(params.Get("after"))
	if err := checkDate(after); err != nil {
		return nil, 0, pagination, httpModel.WrapInvalid(err)
	}

	commits, totalCount, err := a.storeApp.SearchCommit(r.Context(), query, filters, before, after, pagination.PageSize, pagination.Last)
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
		if errors.Is(err, httpModel.ErrInvalid) {
			httperror.BadRequest(w, err)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	var last string
	if len(commits) > 0 {
		last = commits[len(commits)-1].Date.String()
	}

	w.Header().Add("Link", pagination.LinkNextHeader(fmt.Sprintf("%s%s", apiPath, r.URL.Path), r.URL.Query()))
	httpjson.WritePagination(w, http.StatusOK, pagination.PageSize, totalCount, last, commits, httpjson.IsPretty(r))
}

func (a app) handlePostCommits(w http.ResponseWriter, r *http.Request) {
	var commit model.Commit
	if err := httpjson.Parse(r, &commit); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	commit = commit.Sanitize()
	if err := commit.Check(); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if err := a.storeApp.SaveCommit(r.Context(), commit); err != nil {
		httperror.InternalServerError(w, fmt.Errorf("unable to save commit for `%s` with hash `%s`: %s", commit.Repository, commit.Hash, err))
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
