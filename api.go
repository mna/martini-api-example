package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

func GetAlbums(enc Encoder, ar AlbumRepository) string {
	return Must(enc.Encode(toIface(ar.GetAll())...))
}

func GetAlbum(enc Encoder, ar AlbumRepository, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := ar.Get(id)
	if err != nil || al == nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	return 200, Must(enc.Encode(al))
}

func AddAlbum(w http.ResponseWriter, r *http.Request, enc Encoder, ar AlbumRepository) (int, string) {
	al := getPostAlbum(r)
	id, err := ar.Add(al)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the album '%s' from '%s' already exists", al.Title, al.Band))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/albums/%d", id))
		return http.StatusCreated, Must(enc.Encode(al))
	default:
		panic(err)
	}
}

func UpdateAlbum(r *http.Request, enc Encoder, ar AlbumRepository, parms martini.Params) (int, string) {
	al, err := getPutAlbum(r, parms)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	err = ar.Update(al)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the album '%s' from '%s' already exists", al.Title, al.Band))))
	case nil:
		return http.StatusOK, Must(enc.Encode(al))
	default:
		panic(err)
	}
}

func getPostAlbum(r *http.Request) *Album {
	band, title, yrs := r.FormValue("band"), r.FormValue("title"), r.FormValue("year")
	yri, err := strconv.Atoi(yrs)
	if err != nil {
		yri = 0 // Year is optional, set to 0 if invalid/unspecified
	}
	return &Album{
		Band:  band,
		Title: title,
		Year:  yri,
	}
}

func getPutAlbum(r *http.Request, parms martini.Params) (*Album, error) {
	al := getPostAlbum(r)
	id, err := strconv.Atoi(parms["id"])
	if err != nil {
		return nil, err
	}
	al.Id = id
	return al, nil
}

// Martini requires that 2 parameters are returned to treat the first one as the
// status code. Delete is an idempotent action, but this does not mean it should
// always return 204 - No content, idempotence relates to the state of the server
// after the request, not the returned status code. So I return a 404 - Not found
// if the id does not exist.
func DeleteAlbum(enc Encoder, ar AlbumRepository, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := ar.Get(id)
	if err != nil || al == nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	ar.Delete(id)
	return http.StatusNoContent, ""
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
