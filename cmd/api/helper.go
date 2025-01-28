package main

import "encoding/json"

func WrapResponseObject(name string, obj any) ([]byte, error) {
	return json.Marshal(map[string]any{name: obj})
}
