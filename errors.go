package main

import (
	"encoding/xml"
	"fmt"
)

const (
	ErrCodeNotExist      = 1
	ErrCodeAlreadyExists = 2
)

type Error struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Code    int      `json:"code" xml:"code,attr"`
	Message string   `json:"message" xml:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}
