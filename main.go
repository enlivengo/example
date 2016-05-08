package main

import (
	"flag"
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/plugins"
)

func rootHandler(rw http.ResponseWriter, r *http.Request, ev enliven.Enliven) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte("It's working!!"))
}

// Example/Test usage
func main() {
	ev := enliven.New(make(map[string]string))

	// Serving static assets from the ./static/ folder as the /assets/ route
	ev.InitPlugin(plugins.NewStaticAssetPlugin("/assets/", "./static/"))

	// Simple route handler
	ev.AddRoute("/", enliven.RouteHandlerFunc(rootHandler))

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
