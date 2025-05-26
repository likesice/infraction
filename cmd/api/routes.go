package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /health", app.healthHandler())
	mux.Handle("POST /group", app.authMiddleware(app.protectedMiddleware(app.createGroupHandler())))
	mux.Handle("GET /group", app.authMiddleware(app.protectedMiddleware(app.getAllGroupHandler())))
	mux.Handle("POST /login", app.loginHandler())
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	return mux
}
