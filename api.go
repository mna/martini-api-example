package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

// GetAlbums returns the list of albums (possibly filtered).
func GetAlbums(r *http.Request, enc Encoder, db DB) string {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	band, title, yrs := qs.Get("band"), qs.Get("title"), qs.Get("year")
	yri, err := strconv.Atoi(yrs)
	if err != nil {
		// If year is not a valid integer, ignore it
		yri = 0
	}
	if band != "" || title != "" || yri != 0 {
		// At least one filter, use Find()
		return Must(enc.Encode(toIface(db.Find(band, title, yri))...))
	}
	// Otherwise, return all albums
	return Must(enc.Encode(toIface(db.GetAll())...))
}

// GetAlbum returns the requested album.
func GetAlbum(enc Encoder, db DB, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := db.Get(id)
	if err != nil || al == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

// AddAlbum creates the posted album.
func AddAlbum(w http.ResponseWriter, r *http.Request, enc Encoder, db DB) (int, string) {
	al := getPostAlbum(r)
	id, err := db.Add(al)
	switch err {
	case ErrAlreadyExists:
		// Duplicate
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

// UpdateAlbum changes the specified album.
func UpdateAlbum(r *http.Request, enc Encoder, db DB, parms martini.Params) (int, string) {
	al, err := getPutAlbum(r, parms)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	err = db.Update(al)
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

// Parse the request body, load into an Album structure.
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

// Like getPostAlbum, but additionnally, parse and store the `id` query string.
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
func DeleteAlbum(enc Encoder, db DB, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := db.Get(id)
	if err != nil || al == nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
	}
	db.Delete(id)
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
