package main

import (
	"flag"
	"strconv"

	"github.com/hickeroar/enliven"
	_ "github.com/hickeroar/enliven-example/statik"
	"github.com/hickeroar/enliven/apps/admin"
	"github.com/hickeroar/enliven/apps/assets"
	"github.com/hickeroar/enliven/apps/database"
	"github.com/hickeroar/enliven/apps/user"
	"github.com/hickeroar/enliven/middleware/session"
)

func rootHandler(ctx *enliven.Context) {
	val := ctx.Session.Get("increments")

	var value int

	if val == "" {
		value = 1
	} else {
		value, _ = strconv.Atoi(val)
		value++
	}

	newVal := strconv.Itoa(value)
	ctx.Session.Set("increments", newVal)

	tmpl := "{{define \"home\"}}{{template \"header\"}}<div style=\"text-align:center;\">Session Variable: increments = {{.Session.Get \"increments\"}} / {{.Booleans.UserLoggedIn}} / {{.Integers.UserID}}</div>{{template \"footer\"}}{{end}}"
	templates := ctx.Enliven.GetTemplates()
	templates.Parse(tmpl)
	ctx.Template("home")
}

// Example/Test usage
func main() {
	ev := enliven.New(enliven.Config{
		"database_driver":   "postgres",
		"database_host":     "127.0.0.1",
		"database_user":     "postgres",
		"database_dbname":   "enliven",
		"database_password": "postgres",

		"session_redis_address": "127.0.0.1:6379",

		"assets_static_route": "/assets/",
		"assets_static_path":  "./static/",

		"assets_statik_route": "/statik/",
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

	// DATABASE The database app allows you to use....a database
	ev.AddApp(database.NewApp())

	// USER The user app manages the user model/login/session/middleware
	ev.AddApp(user.NewApp())

	// ADMIN The user app manages the admin panel
	ev.AddApp(admin.NewApp())

	// Simple route handler
	ev.AddRoute("/", rootHandler)

	// This is a commented-out example of how you can override the existing header/footer templates.
	// You will most likely want to do this as the built-in one is not meant for general consumption
	// as it has embedded images(base64 encoded)/css/javascript.
	/*
		templates := ev.GetTemplates()
		templates.Parse("{{define \"header\"}}OMG Becky did you see her butt?{{end}}")
		templates.Parse("{{define \"footer\"}}OMG it's so big.{{end}}")
	*/

	port := flag.String("port", ev.GetConfig()["server_port"], "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
