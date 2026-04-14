package static_external

/*
Tried badger here. worked horribly. May be s3/garage could work better.
*/

import (
	"log"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/web/web_utils"
)

func WebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		path = path[1:]

		var content []byte

		var err error

		// code extracting out of third party storage here
		// if err == "not found" {
		// 	http.NotFound(w, r)
		// 	return
		// }

		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		web_utils.MakeStaticFileResp(path, w)

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	log.Println("listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
