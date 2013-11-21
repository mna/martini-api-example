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
	m.Use(RequireHttps)
	m.Use(auth.Basic(AuthToken, ""))
	m.Use(MapEncoder)
	// Setup routes
	r.Get("/albums", GetAlbums)
	// Inject AlbumRepository
	m.MapTo(db, (*AlbumRepository)(nil))
	m.Action(r.Handle)
}

func RequireHttps(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Error(w, "https scheme is required", http.StatusBadRequest)
	}
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
		http.ListenAndServe(":8000", m)
	}()
	http.ListenAndServeTLS(":8001", "cert.pem", "key.pem", m)
}
