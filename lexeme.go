package jajson

type LexemeType uint8

const (
	openCurve LexemeType = iota
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

type lexeme struct {
	typ     LexemeType
	value   []byte
	pos     int
	bytePos int
}
