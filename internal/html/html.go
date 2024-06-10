package html

import (
	"html/template"
	"net/http"

	"github.com/deadshvt/nats-streaming-service/internal/errs"

	"github.com/rs/zerolog"
)

func ParseTemplate(logger zerolog.Logger, w http.ResponseWriter, pattern string, name string, data interface{}) {
	tmpl := template.Must(template.ParseGlob(pattern))
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		msg := errs.WrapError(errs.ErrParseTemplate, err).Error()
		logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
