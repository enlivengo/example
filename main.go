package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/hickeroar/enliven"
	_ "github.com/hickeroar/enliven-example/statik"
	"github.com/hickeroar/enliven/middleware"
	"github.com/hickeroar/enliven/plugins"
	"github.com/jinzhu/gorm"
)

func rootHandler(rw http.ResponseWriter, r *http.Request, ev enliven.Enliven, ctx *enliven.Context) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte("It's working!!"))
}

// User Is a simple user model
type User struct {
	gorm.Model

	Birthday time.Time
	Age      int
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string
}

// Example/Test usage
func main() {
	config := map[string]string{
		"db.driver":   "postgres",
		"db.host":     "127.0.0.1",
		"db.user":     "postgres",
		"db.dbname":   "enliven",
		"db.password": "postgres",
	}

	ev := enliven.New(config)

	// Adding session management middleware
	ev.AddMiddlewareHandler(middleware.NewRedisSessionMiddleware("127.0.0.1:6379", ""))

	// Serving static assets from the ./static/ folder as the /assets/ route
	ev.InitPlugin(plugins.NewStaticAssetPlugin("/assets/", "./static/"))

	// The statik import sets up the data that will be used by the statik filesystem
	// Example: '_ "github.com/hickeroar/enliven-example/statik"'
	ev.InitPlugin(plugins.NewStatikAssetPlugin("/statik/"))

	// Simple route handler
	ev.AddRoute("/", enliven.RouteHandlerFunc(rootHandler))

	ev.GetDatabase().AutoMigrate(&User{})

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
