package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Encoder interface {
	Encode(v ...interface{}) (string, error)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error message as-is.
func MustEncode(data string, err error) string {
	if err != nil {
		panic(err)
	}
	return data
}

type jsonEncoder struct{}

func (_ jsonEncoder) Encode(v ...interface{}) (string, error) {
	var data interface{} = v
	if len(v) == 1 {
		data = v[0]
	}
	b, err := json.Marshal(data)
	return string(b), err
}

type xmlEncoder struct{}

func (_ xmlEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	if _, err := buf.Write([]byte(xml.Header)); err != nil {
		return "", err
	}
	if len(v) > 1 {
		if _, err := buf.Write([]byte("<albums>")); err != nil {
			return "", err
		}
	}
	b, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	if _, err := buf.Write(b); err != nil {
		return "", err
	}
	if len(v) > 1 {
		if _, err := buf.Write([]byte("</albums>")); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

type textEncoder struct{}

func (_ textEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	for _, v := range v {
		if _, err := fmt.Fprintf(&buf, "%s\n", v); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
