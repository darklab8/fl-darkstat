package web

import (
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type Result struct {
	Body       []byte
	StatusCode int
	Resp       *http.Response
}

func Get(client http.Client, url string) (Result, error) {
	res, err := client.Get(url)
	var answer Result
	answer.Resp = res
	answer.StatusCode = res.StatusCode

	if err != nil {
		logus.Log.Error("error making http request: %s\n", typelog.OptError(err))
		return Result{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logus.Log.Error("client: could not read response body: %s\n", typelog.OptError(err))
		return Result{}, err
	}
	answer.Body = resBody
	return answer, nil
}
