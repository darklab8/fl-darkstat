/*
Serve is a very simple static file server in go
Usage:

	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/
package web_static

import (
	"log"
	"net/http"
	"strconv"
)

func WebServer() {
	port := 8000
	directory := "build"

	http.Handle("/", http.FileServer(http.Dir(directory)))

	log.Printf("Serving %s on HTTP port: %d\n", directory, port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
