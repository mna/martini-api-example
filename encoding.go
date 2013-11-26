package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// An Encoder implements an encoding format of values to be sent as response to
// requests on the API endpoints.
type Encoder interface {
	Encode(v ...interface{}) (string, error)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func Must(data string, err error) string {
	if err != nil {
		panic(err)
	}
	return data
}

type jsonEncoder struct{}

// jsonEncoder is an Encoder that produces JSON-formatted responses.
func (_ jsonEncoder) Encode(v ...interface{}) (string, error) {
	var data interface{} = v
	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}
	b, err := json.Marshal(data)
	return string(b), err
}

type xmlEncoder struct{}

// xmlEncoder is an Encoder that produces XML-formatted responses.
func (_ xmlEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	if _, err := buf.Write([]byte(xml.Header)); err != nil {
		return "", err
	}
	if _, err := buf.Write([]byte("<albums>")); err != nil {
		return "", err
	}
	b, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	if _, err := buf.Write(b); err != nil {
		return "", err
	}
	if _, err := buf.Write([]byte("</albums>")); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type textEncoder struct{}

// textEncoder is an Encoder that produces plain text-formatted responses.
func (_ textEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	for _, v := range v {
		if _, err := fmt.Fprintf(&buf, "%s\n", v); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
