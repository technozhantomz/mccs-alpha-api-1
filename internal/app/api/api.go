package api

import (
	"encoding/json"
	"net/http"

	"github.com/technoshantoms/mccs-alpha-api/util/l"
	"go.uber.org/zap"
)

// Respond return an object with specific status as JSON.
func Respond(w http.ResponseWriter, r *http.Request, status int, data ...interface{}) {
	if len(data) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		return
	}

	d := evaluateData(data[0])

	if d == nil {
		w.WriteHeader(status)
		return
	}

	js, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		l.Logger.Error("[ERROR] Marshaling data error:", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func evaluateData(data interface{}) interface{} {
	var value interface{}

	if errors, ok := isMultipleErrors(data); ok {
		value = getErrors(errors)
	} else if e, ok := data.(error); ok {
		value = httpErrors{Errors: []httpError{{Message: e.Error()}}}
	} else {
		value = data
	}

	return value
}

func isMultipleErrors(data interface{}) ([]error, bool) {
	errors, ok := data.([]error)
	if ok {
		return errors, true
	}
	return nil, false
}

func getErrors(data []error) httpErrors {
	errs := httpErrors{}

	for _, err := range data {
		errs.Errors = append(errs.Errors, httpError{Message: err.Error()})
	}

	return errs
}
