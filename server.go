package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/auth"
)

var m *martini.Martini

const AuthToken = "token"

func init() {
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(auth.Basic(AuthToken, ""))
	m.Use(MapEncoder)
	// Setup routes
	r := martini.NewRouter()
	r.Get(`/albums`, GetAlbums)
	r.Get(`/albums/:id`, GetAlbum)
	r.Post(`/albums`, AddAlbum)
	r.Put(`/albums/:id`, UpdateAlbum)
	r.Delete(`/albums/:id`, DeleteAlbum)
	// Inject AlbumRepository
	m.MapTo(db, (*DB)(nil))
	// Add the router action
	m.Action(r.Handle)
}

var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	switch ft {
	case ".xml":
		c.MapTo(xmlEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml")
	case ".text":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	default:
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
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
