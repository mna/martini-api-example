package main

import (
	"encoding/xml"
	"fmt"
	"sync"
)

type AlbumRepository interface {
	Get(id int) *Album
	GetAll() []*Album
	//Find(band, title string, year int) []*Album
	//Add(a *Album) (int, error)
	//Update(a *Album) error
	//Delete(id int) error
}

type albumsDB struct {
	sync.RWMutex
	m map[int]*Album
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
