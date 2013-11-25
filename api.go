package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

func GetAlbums(enc Encoder, ar AlbumRepository) string {
	return MustEncode(enc.Encode(toIface(ar.GetAll())...))
}

func GetAlbum(enc Encoder, ar AlbumRepository, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := ar.Get(id)
	if err != nil || al == nil {
		return http.StatusNotFound, MustEncode(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	return 200, MustEncode(enc.Encode(al))
}

func AddAlbum(w http.ResponseWriter, r *http.Request, enc Encoder, ar AlbumRepository) (int, string) {
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
		return http.StatusConflict, MustEncode(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the album '%s' from '%s' already exists", title, band))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/albums/%d", id))
		return http.StatusCreated, MustEncode(enc.Encode(al))
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
