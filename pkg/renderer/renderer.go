package renderer

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/query"
	"github.com/ViBiOh/httputils/v3/pkg/templates"
)

const (
	faviconPath = "/favicon"
	svgPath     = "/svg"
)

var (
	staticDir = "static"
)

// Input for the renderer
type Input interface {
	GetData(*http.Request) (interface{}, error)
	GetFuncs() template.FuncMap
}

// App of package
type App interface {
	Handler() http.Handler
	IsHandled(*http.Request) bool
}

// Config of package
type Config struct {
	templates *string
}

type app struct {
	tpl     *template.Template
	version string
	input   Input
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		templates: flags.New(prefix, "herodote").Name("Templates").Default("./templates/").Label("HTML Templates folder").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, input Input) (App, error) {
	filesTemplates, err := templates.GetTemplates(strings.TrimSpace(*config.templates), ".html")
	if err != nil {
		return nil, fmt.Errorf("unable to get templates: %s", err)
	}

	return app{
		tpl:     template.Must(template.New("herodote").Funcs(input.GetFuncs()).ParseFiles(filesTemplates...)),
		version: os.Getenv("VERSION"),
		input:   input,
	}, nil
}

func (a app) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, faviconPath) || r.URL.Path == "/robots.txt" || r.URL.Path == "/sitemap.xml" {
			http.ServeFile(w, r, path.Join(staticDir, r.URL.Path))
			return
		}

		if a.tpl != nil && query.IsRoot(r) {
			a.publicHandler(w, r, http.StatusOK, model.ParseMessage(r))
			return
		}

		httperror.NotFound(w)
	})
}

func (a app) IsHandled(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, faviconPath) || r.URL.Path == "/robots.txt" || r.URL.Path == "/sitemap.xml" || query.IsRoot(r)
}
