package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("album already exists")
)

type AlbumRepository interface {
	Get(id int) *Album
	GetAll() []*Album
	//Find(band, title string, year int) []*Album
	Add(a *Album) (int, error)
	//Update(a *Album) error
	Delete(id int)
}

type albumsDB struct {
	sync.RWMutex
	m   map[int]*Album
	seq int
}

var db *albumsDB

func init() {
	db = &albumsDB{
		m: make(map[int]*Album),
	}
	// Fill the database
	db.m[1] = &Album{Id: 1, Band: "Slayer", Title: "Reign In Blood", Year: 1986}
	db.m[2] = &Album{Id: 2, Band: "Slayer", Title: "Seasons In The Abyss", Year: 1990}
	db.m[3] = &Album{Id: 3, Band: "Bruce Springsteen", Title: "Born To Run", Year: 1975}
	db.seq = 3
}

func (db *albumsDB) GetAll() []*Album {
	db.RLock()
	defer db.RUnlock()
	if len(db.m) == 0 {
		return nil
	}
	ar := make([]*Album, len(db.m))
	i := 0
	for _, v := range db.m {
		ar[i] = v
		i++
	}
	return ar
}

func (db *albumsDB) Get(id int) *Album {
	db.RLock()
	defer db.RUnlock()
	return db.m[id]
}

func (db *albumsDB) Add(a *Album) (int, error) {
	db.Lock()
	defer db.Unlock()
	// Return an error if band-title already exists
	if !db.isUnique(a) {
		return 0, ErrAlreadyExists
	}
	// Get the unique ID
	db.seq++
	a.Id = db.seq
	// Store
	db.m[a.Id] = a
	return a.Id, nil
}

func (db *albumsDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()
	delete(db.m, id)
}

func (db *albumsDB) isUnique(a *Album) bool {
	for _, v := range db.m {
		if v.Band == a.Band && v.Title == a.Title {
			return false
		}
	}
	return true
}

type Album struct {
	XMLName xml.Name `json:"-" xml:"album"`
	Id      int      `json:"id" xml:"id,attr"`
	Band    string   `json:"band" xml:"band"`
	Title   string   `json:"title" xml:"title"`
	Year    int      `json:"year" xml:"year"`
}

func (a *Album) String() string {
	return fmt.Sprintf("%s - %s (%d)", a.Band, a.Title, a.Year)
}
