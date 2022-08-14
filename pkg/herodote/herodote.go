package herodote

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v4/pkg/httperror"
	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	httpModel "github.com/ViBiOh/httputils/v4/pkg/model"
	"github.com/ViBiOh/httputils/v4/pkg/query"
	"github.com/ViBiOh/httputils/v4/pkg/renderer"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

const (
	isoDateLayout = "2006-01-02"

	apiPath     = "/api"
	commitsPath = "/commits"
)

// ErrAuthentificationFailed occurs when secret is invalid
var ErrAuthentificationFailed = errors.New("invalid secret provided")

// App of package
type App struct {
	tracer     trace.Tracer
	apiHandler http.Handler
	colors     map[string]string
	storeApp   store.App
	secret     string
}

// Config of package
type Config struct {
	secret *string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		secret: flags.String(fs, prefix, "herodote", "HttpSecret", "HTTP Secret Key for Update", "", nil),
	}
}

// New creates new App from Config
func New(config Config, storeApp store.App, tracer trace.Tracer) (App, error) {
	if len(*config.secret) == 0 {
		return App{}, errors.New("http secret is required")
	}

	if !storeApp.Enabled() {
		return App{}, errors.New("store is required")
	}

	app := App{
		secret:   *config.secret,
		storeApp: storeApp,
		tracer:   tracer,
		colors:   make(map[string]string),
	}

	app.apiHandler = http.StripPrefix(apiPath, app.Handler())

	return app, nil
}

// Handler for request. Should be use with net/http
func (a App) Handler() http.Handler {
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

// TemplateFunc used for rendering UI
func (a App) TemplateFunc(w http.ResponseWriter, r *http.Request) (renderer.Page, error) {
	if strings.HasPrefix(r.URL.Path, apiPath) {
		a.apiHandler.ServeHTTP(w, r)
		return renderer.Page{}, nil
	}

	commits, _, _, err := a.listCommits(r)
	if err != nil {
		return renderer.NewPage("", http.StatusInternalServerError, nil), err
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return renderer.NewPage("", http.StatusInternalServerError, nil), fmt.Errorf("parse query: %s", err)
	}

	filters, err := a.storeApp.ListFilters(r.Context())
	if err != nil {
		return renderer.NewPage("", http.StatusInternalServerError, nil), fmt.Errorf("list filters: %s", err)
	}

	return renderer.NewPage("public", http.StatusOK, map[string]any{
		"Path":         r.URL.Path,
		"Filters":      params,
		"Repositories": filters["repository"],
		"Types":        filters["type"],
		"Components":   filters["component"],
		"Colors":       repositoriesColors,
		"Commits":      commits,
		"Now":          time.Now(),
	}), nil
}

func (a App) listCommits(r *http.Request) ([]model.Commit, uint, query.Pagination, error) {
	ctx, end := tracer.StartSpan(r.Context(), a.tracer, "list commits", trace.WithSpanKind(trace.SpanKindInternal))
	defer end()

	pagination, err := query.ParsePagination(r, 50, 100)
	if err != nil {
		return nil, 0, pagination, httpModel.WrapInvalid(err)
	}

	params := r.URL.Query()

	searchQuery := strings.TrimSpace(params.Get("q"))
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

	commits, totalCount, err := a.storeApp.SearchCommit(ctx, searchQuery, filters, before, after, pagination.PageSize, pagination.Last)
	return commits, totalCount, pagination, err
}

func (a App) handleCommits(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.handlePostCommits(w, r)
	} else if r.Method == http.MethodGet {
		a.handleGetCommits(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a App) handleGetCommits(w http.ResponseWriter, r *http.Request) {
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
	httpjson.WritePagination(w, http.StatusOK, pagination.PageSize, totalCount, last, commits)
}

func (a App) handlePostCommits(w http.ResponseWriter, r *http.Request) {
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
		httperror.InternalServerError(w, fmt.Errorf("save commit for `%s` with hash `%s`: %s", commit.Repository, commit.Hash, err))
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
		return fmt.Errorf("parse date: %s", err)
	}

	return nil
}
