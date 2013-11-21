package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Encoder interface {
	Encode(w http.ResponseWriter, v interface{}) error
}

type jsonEncoder struct{}

func (_ jsonEncoder) Encode(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

type xmlEncoder struct{}

func (_ xmlEncoder) Encode(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	return enc.Encode(v)
}

type textEncoder struct{}

func (_ textEncoder) Encode(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(fmt.Sprintf("%s", v)))
	return err
}
