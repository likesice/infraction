package main

import (
	"fmt"
	"infraction.mageis.net/internal/store"
	"net/http"
	"strings"
)

func (app *Application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteHttpError(w, http.StatusUnauthorized, fmt.Errorf(""))
			return
		}
		// Extract the actual authentication token from the header parts.
		token := headerParts[1]

		user, err := app.registry.UserStore.GetForToken(token)
		if err != nil {
			WriteHttpError(w, http.StatusUnauthorized, fmt.Errorf(""))
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) protectedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			WriteHttpError(w, http.StatusUnauthorized, fmt.Errorf(""))
			return
		}
		next.ServeHTTP(w, r)
	})
}
