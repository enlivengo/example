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
	ev := enliven.New(map[string]string{
		"db.driver":   "postgres",
		"db.host":     "127.0.0.1",
		"db.user":     "postgres",
		"db.dbname":   "enliven",
		"db.password": "postgres",

		"session.redis.address":  "127.0.0.1:6379",
		"session.redis.password": "",

		"static.assets.route": "/assets/",
		"static.assets.path":  "./static/",

		"statik.assets.route": "/statik/",
	})

	// Adding session management middleware
	ev.AddMiddleware(middleware.NewRedisSessionMiddleware(ev.GetConfig()))

	// Serving static assets from the ./static/ folder as the /assets/ route
	ev.InitPlugin(plugins.NewStaticAssetsPlugin(ev.GetConfig()))

	// The statik import sets up the data that will be used by the statik filesystem. Read Statik documentation
	ev.InitPlugin(plugins.NewStatikAssetsPlugin(ev.GetConfig()))

	// Simple route handler
	ev.AddRoute("/", rootHandler)

	// We can use gorm to automigrate our database on application start
	ev.GetDatabase().AutoMigrate(&User{})

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
