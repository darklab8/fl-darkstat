package apiutils

import (
	"encoding/json"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
)

func ReturnJson(w *http.ResponseWriter, data any) {
	(*w).Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(*w).Encode(data)
	if logus.Log.CheckError(err, "should be marshable") {
		json.NewEncoder(*w).Encode(struct {
			Error string
		}{
			Error: "not marshable for some reason",
		})
		(*w).WriteHeader(http.StatusInternalServerError)
	}
}
