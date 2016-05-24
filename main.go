package main

import (
	"strconv"

	"github.com/enlivengo/admin"
	"github.com/enlivengo/database"
	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/apps/statik"
	"github.com/enlivengo/enliven/config"
	"github.com/enlivengo/enliven/middleware/session"
	_ "github.com/enlivengo/example/statik"
	"github.com/enlivengo/static"
	"github.com/enlivengo/user"
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

	ctx.Template("home")
}

// Example/Test usage
func main() {
	ev := enliven.New(config.Config{
		"email_smtp_host":    "localhost",
		"email_smtp_auth":    "none",
		"email_from_default": "noreply@enliven.app",

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
	ev.AddMiddleware(session.NewRedisStorageMiddleware())
	//ev.AddMiddleware(session.NewFileStorageMiddleware())
	//ev.AddMiddleware(session.NewMemoryStorageMiddleware())

	// Serving static assets from the ./static/ folder at the /assets/ route
	ev.AddApp(static.NewApp())

	// The statik import at the top of this file sets up the data that will be used by the statik filesystem.
	// Read Statik documentation
	ev.AddApp(statik.NewApp())

	// DATABASE The database app allows you to use....a database
	ev.AddApp(database.NewApp())

	// USER The user app manages the user model/login/session/middleware
	ev.AddApp(user.NewApp())

	// ADMIN The admin app manages the admin panel
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

	ev.Run()
}
