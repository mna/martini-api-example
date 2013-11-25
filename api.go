package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error message as-is.

// TODO : Return errors in same format as requested (use Error struct and Encoder)
func GetAlbums(enc Encoder, ar AlbumRepository) string {
	data, err := enc.Encode(toIface(ar.GetAll())...)
	if err != nil {
		panic(err)
	}
	return data
}

func GetAlbum(enc Encoder, ar AlbumRepository, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		return http.StatusNotFound, "album not found"
	}
	al := ar.Get(id)
	if al == nil {
		return http.StatusNotFound, "album not found"
	}
	data, err := enc.Encode(al)
	if err != nil {
		panic(err)
	}
	return 200, data
}

func AddAlbum(w http.ResponseWriter, r *http.Request, enc Encoder, ar AlbumRepository) {
	band, title, yrs := r.FormValue("band"), r.FormValue("title"), r.FormValue("year")
	yri, err := strconv.Atoi(yrs)
	if err != nil {
		yri = 0 // Year is optional, set to 0 if invalid/unspecified
	}
	al := &Album{
		Band:  band,
		Title: title,
		Year:  yri,
	}
	id, err := ar.Add(al)
	switch err {
	case ErrAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		data, err := enc.Encode(al)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Location", fmt.Sprintf("/albums/%d", id))
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(data)); err != nil {
			panic(err)
		}
	default:
		panic(err)
	}
}

func toIface(v []*Album) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}
