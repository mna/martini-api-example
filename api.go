package main

import (
	"net/http"
)

func GetAlbums(w http.ResponseWriter, enc Encoder, ar AlbumRepository) {
	if err := enc.Encode(w, ar.GetAll()); err != nil {
		// TODO : You probably don't want to expose internal errors like this
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
