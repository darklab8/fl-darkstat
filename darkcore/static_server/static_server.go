/*
it was addded for testing idea in optimizing Discovery darkstat running through aggressive caching to filesystem
it did not work well, because that much writing to filesystem (of 360'000 files with 4Gb total sum) is expensive
and makes entire system lagging. I leave it behind for inspiration

This experience lead me to believe there are two possible way to fix Discovery darkstat running for zero downtime updates
1) Improve in memory refresh for `darkstat web` running
2) try static assets cron idea again but utilize Badger as your database of static assets!
*/
package static_server

import (
	"log"
	"net/http"
	"strconv"
)

func StaticServer() {
	port := 8000
	directory := "build"

	fileServer := http.FileServer(http.Dir(directory))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, directory+"/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("Serving %s on HTTP port: %d\n", directory, port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
