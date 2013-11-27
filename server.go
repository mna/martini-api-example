package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/auth"
)

// The one and only access token! In real-life scenarios, a more complex authentication
// middleware than auth.Basic should be used, obviously.
const AuthToken = "token"

// The one and only martini instance.
var m *martini.Martini

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
	// Inject database
	m.MapTo(db, (*DB)(nil))
	// Add the router action
	m.Action(r.Handle)
}

// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	// Inject the requested encoder
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
		// scalability purposes, so it's not like it will block otherwise valid routes.
		//
		// It is also common practice to use a different subdomain so that cookies are
		// not transfered with every API request.
		// So with that in mind, it seems reasonable to refuse each and every request
		// on the non-https server, regardless of the route. This could of course be done
		// on a reverse-proxy in front of this web server.
		//
		if err := http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "https scheme is required", http.StatusBadRequest)
		})); err != nil {
			log.Fatal(err)
		}
	}()

	// Listen on https: with the preconfigured martini instance. The certificate files
	// can be created using this command in this repository's root directory:
	//
	// go run /path/to/goroot/src/pkg/crypto/tls/generate_cert.go --host="localhost"
	//
	if err := http.ListenAndServeTLS(":8001", "cert.pem", "key.pem", m); err != nil {
		log.Fatal(err)
	}
}
