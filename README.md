# A RESTful API example with Martini

## Install and run

To install this example application, run the usual:

`go get github.com/PuerkitoBio/martini-api-example`

To run locally, it requires `cert.pem` and `key.pem` files (for TLS/https). You can generate the pair of files using this simple command, adapted to your environment for the GOROOT part:

`go run /path/to/goroot/src/pkg/crypto/tls/generate_cert.go --host="localhost"`

When running the API server (which will be named `martini-api-example`, in your $GOPATH/bin directory), it will look for the `*.pem` files in the current directory. You can adapt the code as you see fit, if you want to specify, for example, an absolute path (it is in `server.go`).

## Blog post

This repository is the companion source code for my [blog post on 0value.com][blog]. It shows how to build a RESTful API using the [martini package][martini]. Please note that martini is currently a moving target, so please file an issue (or submit a PR) if the code doesn't work anymore with more recent versions of martini.

## License

The [BDS 3-clause license][bsd].

[martini]: https://github.com/codegangsta/martini
[blog]: http://0value.com/build-a-restful-API-with-Martini
[bsd]: http://opensource.org/licenses/BSD-3-Clause
