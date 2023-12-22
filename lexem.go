package jajson

type LexemType uint8

const (
	openCurve LexemType = iota
	closeCurve
	openBracket
	closeBracket
	colon
	comma
	String
	Int
	Float
	Bool

	Object
	Array
	Err
)

type lexem struct {
	typ     LexemType
	value   []byte
	pos     int
	bytePos int
}
