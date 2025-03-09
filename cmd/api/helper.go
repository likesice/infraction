package main

import (
	"encoding/json"
	"net/http"
)

func WrapResponseObject(name string, obj any) ([]byte, error) {
	return json.Marshal(map[string]any{name: obj})
}

func WriteHttpError(w http.ResponseWriter, code int, err error) {
	res, err := WrapResponseObject("error", err.Error())
	if err != nil {
		//TODO: write generic err obj
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
}
