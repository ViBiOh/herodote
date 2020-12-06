package renderer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ViBiOh/herodote/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/templates"
)

func redirectWithMessage(w http.ResponseWriter, r *http.Request, path, message string) {
	http.Redirect(w, r, fmt.Sprintf("%s?messageContent=%s", path, url.QueryEscape(message)), http.StatusFound)
}

func (a app) publicHandler(w http.ResponseWriter, r *http.Request, status int, message model.Message) {
	datas, err := a.input.GetData(r)
	if err != nil {
		a.errorHandler(w, http.StatusBadRequest, err)
	}

	content := map[string]interface{}{
		"Version": a.version,
		"Data":    datas,
	}

	if len(message.Content) > 0 {
		content["Message"] = message
	}

	if err := templates.ResponseHTMLTemplate(a.tpl.Lookup("public"), w, content, status); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func (a app) errorHandler(w http.ResponseWriter, status int, err error) {
	logger.Error("%s", err)

	content := map[string]interface{}{
		"Version": a.version,
	}

	if err != nil {
		message := err.Error()
		subMessages := ""

		if errors.Is(err, model.ErrInvalid) {
			status = http.StatusBadRequest
		} else if errors.Is(err, model.ErrInternalError) {
			status = http.StatusInternalServerError
			message = "Oops! Something went wrong."
		}

		content["Message"] = model.NewErrorMessage(message)
		if len(subMessages) > 0 {
			content["Errors"] = strings.Split(subMessages, ", ")
		}
	}

	if err := templates.ResponseHTMLTemplate(a.tpl.Lookup("error"), w, content, status); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func (a app) svg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl := a.tpl.Lookup(fmt.Sprintf("svg-%s", strings.Trim(r.URL.Path, "/")))
		if tpl == nil {
			httperror.NotFound(w)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		if err := templates.WriteTemplate(tpl, w, r.URL.Query().Get("fill"), "text/xml"); err != nil {
			httperror.InternalServerError(w, err)
		}
	})
}
