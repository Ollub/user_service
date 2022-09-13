package http_utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func FromBody[T interface{}](r *http.Request) (*T, error) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	out := new(T)
	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
