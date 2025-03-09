package main

import (
	"context"
	"net/http"
)

type UserKey string

const UserCtxKey = UserKey("user")

func (app *Application) dummyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := app.registry.UserStore.GetUser("dummy")
		r = r.WithContext(context.WithValue(r.Context(), UserCtxKey, user))
		next.ServeHTTP(w, r)
	})
}
