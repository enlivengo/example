package main

import (
	"flag"
	"strconv"

	"github.com/hickeroar/enliven"
	_ "github.com/hickeroar/enliven-example/statik"
	"github.com/hickeroar/enliven/apps/assets"
	"github.com/hickeroar/enliven/apps/user"
	"github.com/hickeroar/enliven/middleware/session"
)

func rootHandler(ctx *enliven.Context) {
	val := ctx.Session.Get("increments")

	var value int

	if val == "" {
		val = "1"
		value = 1
	} else {
		value, _ = strconv.Atoi(val)
		value++
	}

	newVal := strconv.Itoa(value)
	ctx.Session.Set("increments", newVal)

	ctx.String("Session Variable: increments = " + val + " / " + ctx.Items["UserLoggedIn"] + " / " + ctx.Items["UserID"])
}

// Example/Test usage
func main() {
	ev := enliven.New(enliven.Config{
		"db.driver":   "postgres",
		"db.host":     "127.0.0.1",
		"db.user":     "postgres",
		"db.dbname":   "enliven",
		"db.password": "postgres",

		"session.redis.address": "127.0.0.1:6379",

		"assets.static.route": "/assets/",
		"assets.static.path":  "./static/",

		"assets.statik.route": "/statik/",
	})

	// Adding session management middleware
	//ev.AddMiddleware(session.NewRedisStorageMiddleware(ev.GetConfig()))
	//ev.AddMiddleware(session.NewFileStorageMiddleware(ev.GetConfig()))
	ev.AddMiddleware(session.NewMemoryStorageMiddleware(ev.GetConfig()))

	// Serving static assets from the ./static/ folder as the /assets/ route
	ev.AddApp(assets.NewStaticApp())

	// The statik import at the top of this file sets up the data that will be used by the statik filesystem.
	// Read Statik documentation
	ev.AddApp(assets.NewStatikApp())

	// The user app manages the user model/login/session/middleware
	ev.AddApp(user.NewApp())

	// Simple route handler
	ev.AddRoute("/", rootHandler)

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
