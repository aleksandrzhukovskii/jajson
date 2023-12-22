package jajson

import (
	"errors"
	"fmt"
)

type Error struct {
	pos int
	err error
}

func (e Error) New(pos int) Error {
	e.pos = pos
	return e
}

func (e Error) Error() string {
	return fmt.Sprintf("Pos: %d. Error: %s", e.pos, e.err.Error())
}

var ErrorUnexpected = Error{err: errors.New("unexpected symbol or end of JSON")}
var ErrorRune = Error{err: errors.New("cannot parse next rune")}
var ErrorWrongQuote = Error{err: errors.New("wrong quotation")}
var ErrorEmptyJSON = Error{err: errors.New("JSON is empty")}
