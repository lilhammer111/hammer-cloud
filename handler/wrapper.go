package handler

import (
	"fmt"
	"net/http"
)

func ParseRequestParams(r *http.Request, keys ...string) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for _, key := range keys {
		value := r.Form.Get(key)
		if value == "" {
			return nil, fmt.Errorf("missing param: %s", key)
		}
		params[key] = value
	}
	return params, nil
}
