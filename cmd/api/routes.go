package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /health", app.healthHandler())
	mux.Handle("POST /group", app.dummyAuthMiddleware(app.createGroupHandler()))
	mux.Handle("GET /group", app.dummyAuthMiddleware(app.getAllGroupHandler()))
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	return mux
}
