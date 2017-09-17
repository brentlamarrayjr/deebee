package deebee

import (
	"fmt"
)

//Error is the type of error thrown by package
type Error struct {
	number      int
	description string
}

func (e Error) Error() string {
	return fmt.Sprintf("EscuelleError(%d): %s", e.number, e.description)
}

//ErrInvalidDSN is thrown when provided DSN is invalid
var ErrInvalidDSN = Error{number: 0, description: "Invalid DSN."}

//ErrTypeNotSupported is thrown when reflect.Type is not supported by function
var ErrTypeNotSupported = Error{number: 1, description: "Type not supported."}

//ErrMethodNotImplemented is thrown when reflect.Type is not supported by function
var ErrMethodNotImplemented = Error{number: 2, description: "Method not supported."}
