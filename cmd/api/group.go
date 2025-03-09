package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"infraction.mageis.net/internal/store"
	"net/http"
)

var (
	ErrGroupValidation = fmt.Errorf("group object validation failed")
	ErrGroupNoUser     = fmt.Errorf("group request has no valid user")
)

func (app *Application) createGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var group store.Group
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			WriteHttpError(w, http.StatusBadRequest, ErrGroupValidation)
			app.logger.Error("createGroupHandler: ", errors.Join(ErrGroupValidation, err))
			return
		}
		u, ok := r.Context().Value(UserCtxKey).(*store.User)
		if !ok {
			WriteHttpError(w, http.StatusForbidden, ErrGroupNoUser)
			app.logger.Error("createGroupHandler: ", ErrGroupNoUser)
			return
		}
		group.Members = []int64{u.Id}
		err = app.registry.GroupStore.Insert(&group)
		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError, err)
			app.logger.Error("createGroupHandler: ", err)
			return
		}
		responseObject, err := WrapResponseObject("result", group)
		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError, err)
			app.logger.Error("createGroupHandler: ", err)
			return
		}
		w.Write(responseObject)
	}
}

func (app *Application) getAllGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(UserCtxKey).(*store.User)
		if !ok {
			WriteHttpError(w, http.StatusForbidden, ErrGroupNoUser)
			app.logger.Error("getAllGroupHandler: ", ErrGroupNoUser)
			return
		}
		groups, err := app.registry.GroupStore.SelectAll(u)
		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError, err)
			app.logger.Error("getAllGroupHandler: ", err)
			return
		}

		responseObject, err := WrapResponseObject("result", groups)

		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError, err)
			app.logger.Error("getAllGroupHandler: ", err)
			return
		}
		w.Write(responseObject)
	}
}
