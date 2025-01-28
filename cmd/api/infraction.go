package main

import (
	"encoding/json"
	"infraction.mageis.net/internal/data"
	"net/http"
)

func (app *Application) createInfractionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var infraction data.Infraction
		err := json.NewDecoder(r.Body).Decode(&infraction)
		if err != nil {
			app.logger.Error("WTF: ", err)
		}
		u, ok := r.Context().Value(UserCtxKey).(*data.User)
		if !ok {
			app.logger.Error("FUCK: ", err)
		}
		infraction.User = u.Id
		err = app.registry.InfractionRepository.Insert(&infraction)
		if err != nil {
			http.Error(w, "oops", 500)
			app.logger.Error("FUCK: ", err)
			return
		}
	}
}

func (app *Application) getAllInfractionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(UserCtxKey).(*data.User)
		if !ok {
			app.logger.Error("FUCK: ")
		}
		infractions, err := app.registry.InfractionRepository.SelectAll(u)
		if err != nil {
			return
		}

		responseObject, err := WrapResponseObject("result", infractions)
		if err != nil {
			http.Error(w, "oops", 500)
			app.logger.Error("FUCK: ", err)
			return
		}
		w.Write(responseObject)
	}
}
