package backend

import (
	"encoding/json"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/settings/types"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type Backend struct {
}

func NewBackend() *Backend {
	return &Backend{}
}

type ErrorMessage struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

func NewErrorMsg(err error) string {
	result, err := json.Marshal(&ErrorMessage{
		Msg:  err.Error(),
		Type: "scrappy_data_serialization_error",
	})
	logus.Log.CheckError(err, "failed to marshal error")
	return string(result)
}

func (app *Backend) RegisterBack() {
	endpoint_ping := NewEndpointPing(app)
	http.HandleFunc(string(endpoint_ping.Url), endpoint_ping.Handler)
}

type Endpoint struct {
	Url     types.Url
	Handler func(w http.ResponseWriter, r *http.Request)
}
