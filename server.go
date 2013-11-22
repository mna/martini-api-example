package main

import (
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/auth"
)

var m *martini.Martini

const AuthToken = "token"

func init() {
	m = martini.New()
	r := martini.NewRouter()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(auth.Basic(AuthToken, ""))
	m.Use(MapEncoder)
	// Setup routes
	// TODO : Support extension-style format (.json, etc.)
	r.Get("/albums", GetAlbums)
	// Inject AlbumRepository
	m.MapTo(db, (*AlbumRepository)(nil))
	m.Action(r.Handle)
}

func MapEncoder(c martini.Context, r *http.Request) {
	vals := r.URL.Query()
	switch vals.Get("f") {
	case "xml":
		c.MapTo(xmlEncoder{}, (*Encoder)(nil))
	case "text":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
	default:
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
	}
}

func main() {
	go func() {
		// Listen on http: to raise an error and indicate that https: is required.
		//
		// This could also be achieved by passing the same `m` martini instance as
		// used by the https server, and by using a middleware that checks for https
		// and returns an error if it is not a secure connection. This would have the benefit
		// of handling only the defined routes. However, it is common practice to define
		// APIs on separate web servers from the web (html) pages, for maintenance and
		// scalability purposes, among other things. It is also common practice to use
		// a different subdomain so that cookies are not transfered with every API request.
		// So with that in mind, it seems reasonable to refuse each and every request
		// on the non-https server, regardless of the route.
		//
		http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "https scheme is required", http.StatusBadRequest)
		}))
	}()

	// Listen on https: with the preconfigured martini instance.
	http.ListenAndServeTLS(":8001", "cert.pem", "key.pem", m)
}
