package apiutils

import (
	"encoding/json"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
)

func ReturnJson(w *http.ResponseWriter, data any) {
	(*w).Header().Set("Content-Type", "application/json")

	marshaled, err := json.Marshal(data)
	if logus.Log.CheckError(err, "should be marshable") {
		err := json.NewEncoder(*w).Encode(struct {
			Error string
		}{
			Error: "not marshable for some reason",
		})
		logus.Log.CheckWarn(err, "failed to encode error response in return json")
		(*w).WriteHeader(http.StatusInternalServerError)
	}

	(*w).Write(marshaled)
}
