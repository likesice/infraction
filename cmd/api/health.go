package main

import (
	"context"
	"fmt"
	"net/http"
)

func (app *Application) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := WrapResponseObject("data", map[string]string{
			"status":  "OK",
			"version": "1.0.0-alpha",
		})
		if err != nil {
			//TODO: unify error handling/logging
			err = fmt.Errorf("could not wrap response object: %w", err)
			r.WithContext(context.WithValue(r.Context(), "err", err))
		}

		_, err = w.Write(buf)
		if err != nil {
			err = fmt.Errorf("could not serialize wrapped response object: %w", err)
			r.WithContext(context.WithValue(r.Context(), "err", err))
		}
	}
}
