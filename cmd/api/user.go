package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrUserAuthentication = fmt.Errorf("user authentication failed")

func (app *Application) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDto struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		//TODO: validation
		err := json.NewDecoder(r.Body).Decode(&loginDto)
		if err != nil {
			WriteHttpError(w, http.StatusForbidden, ErrUserAuthentication)
			app.logger.Error("loginHandler:", errors.Join(ErrUserAuthentication, err))
			return
		}

		user, err := app.registry.UserStore.GetUser(loginDto.Email)

		if err != nil {
			//TODO: adjust err
			WriteHttpError(w, http.StatusForbidden, ErrUserAuthentication)
			app.logger.Error("loginHandler:", errors.Join(ErrUserAuthentication, err))
			return
		}

		matches, err := user.Password.Matches(loginDto.Password)
		if err != nil {
			//TODO: adjust err
			WriteHttpError(w, http.StatusForbidden, ErrUserAuthentication)
			app.logger.Error("loginHandler:", errors.Join(ErrUserAuthentication, err))
			return
		}
		if !matches {
			//TODO: adjust err
			WriteHttpError(w, http.StatusForbidden, ErrUserAuthentication)
			app.logger.Error("loginHandler:", errors.Join(ErrUserAuthentication, err))
			return
		}

		session, err := app.registry.SessionStore.New(user)
		if err != nil {
			//TODO: adjust err
			WriteHttpError(w, http.StatusForbidden, ErrUserAuthentication)
			app.logger.Error("loginHandler:", errors.Join(ErrUserAuthentication, err))
			return
		}

		responseObject, err := WrapResponseObject("data", session)

		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError, err)
			app.logger.Error("loginHandler: ", err)
			return
		}
		w.Write(responseObject)
	}
}
