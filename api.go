package main

import (
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

func GetAlbums(w http.ResponseWriter, enc Encoder, ar AlbumRepository) {
	if err := enc.Encode(w, toIface(ar.GetAll())...); err != nil {
		// TODO : You probably don't want to expose internal errors like this
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetAlbum(w http.ResponseWriter, enc Encoder, ar AlbumRepository, parms martini.Params) {
	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		http.Error(w, "404 album not found", http.StatusNotFound)
		return
	}
	al := ar.Get(id)
	if al == nil {
		http.Error(w, "404 album not found", http.StatusNotFound)
		return
	}
	if err = enc.Encode(w, al); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
