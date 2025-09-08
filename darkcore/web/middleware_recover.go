package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-utils/typelog"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				logus.Log.Error(
					"web middleware recovery caught smth",
					typelog.Any("err", fmt.Sprint(err)),
					typelog.Any("stack", string(debug.Stack())),
				)
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error. Check logs for stacktraces",
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}

		}()

		next.ServeHTTP(w, r)

	})
}
