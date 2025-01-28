package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /health", app.healthHandler())
	mux.Handle("POST /infraction", app.dummyAuthMiddleware(app.createInfractionHandler()))
	mux.Handle("GET /infraction", app.dummyAuthMiddleware(app.getAllInfractionHandler()))
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("oopsie, nothing to see here 💩"))
	}))
	return mux
}
