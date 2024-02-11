package front_utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type Request[T any] struct {
	url    string
	method string
	body   []byte
}

func NewRequest[T any](url string, method string, opts ...RequestParam[T]) *Request[T] {
	r := &Request[T]{
		url:    url,
		method: method,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type RequestParam[T any] func(r *Request[T])

func WithBody[T any](body []byte) RequestParam[T] {
	return func(r *Request[T]) {
		r.body = body
	}
}

func RequestRun[T any](r *Request[T]) (*T, error) {
	var result T

	var res *http.Response
	var err error
	switch r.method {
	case http.MethodGet:
		res, err = http.Get(r.url)
	case http.MethodPost:
		fmt.Println("TODO")
		res, err = http.Post(
			r.url,
			"application/json",
			bytes.NewReader(r.body))
	default:
		return nil, errors.New("not supported method")
	}

	if logus.Log.CheckError(err, "failed to request vars.json") {
		return &result, err
	}

	resBody, err := io.ReadAll(res.Body)
	if logus.Log.CheckError(err, "failed to read hook env vars body") {
		return &result, err
	}

	err = json.Unmarshal(resBody, &result)
	logus.Log.CheckError(err, "failed to parse environment variables at front")
	return &result, nil
}
