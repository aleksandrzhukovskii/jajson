package jajson

import (
	"unicode"
	"unicode/utf8"
)

var rue = []rune("rue")
var alse = []rune("alse")
var runeToType = map[rune]LexemType{'{': openCurve, '}': closeCurve, '[': openBracket, ']': closeBracket, ':': colon, ',': comma}

type lexer struct {
	data    []byte
	pos     int
	bytePos int
}

func newLexer(data []byte) *lexer {
	return &lexer{
		data:    data,
		pos:     0,
		bytePos: 0,
	}
}

func (t *lexer) nextToken() (lexem, error) {
	if len(t.data) == 0 {
		return lexem{}, ErrorUnexpected.New(t.pos)
	}
	r, size := utf8.DecodeRune(t.data)
	if r == utf8.RuneError {
		return lexem{}, ErrorRune.New(t.pos)
	}
	before := t.data
	t.data = t.data[size:]
	for unicode.IsSpace(r) && len(t.data) > 0 {
		t.pos++
		t.bytePos += size
		r, size = utf8.DecodeRune(t.data)
		if r == utf8.RuneError {
			return lexem{}, ErrorRune.New(t.pos)
		}
		before = t.data
		t.data = t.data[size:]
	}

	switch r {
	case '{', '}', '[', ']', ':', ',':
		defer func() { t.pos++; t.bytePos += size }()
		return lexem{typ: runeToType[r], pos: t.pos, bytePos: t.bytePos}, nil
	case 't':
		if err := t.skip(rue); err != nil {
			return lexem{}, err
		}
		byteLen := len(before) - len(t.data)
		defer func() { t.pos += 4; t.bytePos += byteLen }()
		return lexem{typ: Bool, pos: t.pos, value: before[:byteLen], bytePos: t.bytePos}, nil
	case 'f':
		if err := t.skip(alse); err != nil {
			return lexem{}, err
		}
		byteLen := len(before) - len(t.data)
		defer func() { t.pos += 5; t.bytePos += byteLen }()
		return lexem{typ: Bool, pos: t.pos, value: before[:byteLen], bytePos: t.bytePos}, nil
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		ret, float, err := t.skipNum(r == '-')
		if err != nil {
			return lexem{}, err
		}
		byteLen := len(before) - len(t.data)
		defer func() { t.pos += ret + 1; t.bytePos += byteLen }()
		var typ LexemType
		if float {
			typ = Float
		} else {
			typ = Int
		}
		return lexem{typ: typ, pos: t.pos, value: before[:byteLen], bytePos: t.bytePos}, nil
	case '"':
		ret, err := t.skipString()
		if err != nil {
			return lexem{}, err
		}
		byteLen := len(before) - len(t.data)
		defer func() { t.pos += ret + 2; t.bytePos += byteLen }()
		return lexem{typ: String, pos: t.pos, value: before[:byteLen], bytePos: t.bytePos}, nil
	default:
		return lexem{}, ErrorUnexpected.New(t.pos)
	}
}

func (t *lexer) skip(str []rune) error {
	for i := 0; i < len(str); i++ {
		if len(t.data) == 0 {
			return ErrorUnexpected.New(t.pos)
		}
		r, size := utf8.DecodeRune(t.data)
		if r == utf8.RuneError {
			return ErrorRune.New(t.pos)
		}
		t.data = t.data[size:]
		if r != str[i] {
			return ErrorUnexpected.New(t.pos)
		}
	}
	return nil
}

func (t *lexer) skipNum(zeroCritical bool) (int, bool, error) {
	ret := 0
	float := false
	point := false
	for {
		if len(t.data) == 0 {
			if ret == 0 && zeroCritical {
				return 0, false, ErrorUnexpected.New(t.pos)
			} else {
				return ret, float, nil
			}
		}
		r, size := utf8.DecodeRune(t.data)
		if r == utf8.RuneError {
			if ret == 0 && zeroCritical {
				return 0, false, ErrorRune.New(t.pos)
			} else if point {
				return 0, false, ErrorUnexpected.New(t.pos)
			} else {
				return ret, float, nil
			}
		}
		if r != '.' && r != ',' && !(r >= '0' && r <= '9') {
			if point {
				return 0, false, ErrorUnexpected.New(t.pos)
			}
			if ret == 0 && zeroCritical {
				return 0, false, ErrorUnexpected.New(t.pos)
			} else {
				return ret, float, nil
			}
		}
		point = false
		if r == '.' || r == ',' {
			point = true
			if float {
				return 0, false, ErrorUnexpected.New(t.pos)
			} else {
				float = true
			}
		}
		t.data = t.data[size:]
		ret++
	}
}

func (t *lexer) skipString() (int, error) {
	if len(t.data) < 1 {
		return 0, ErrorUnexpected.New(t.pos)
	}
	ret := 0
	for len(t.data) > 0 && t.data[0] != '"' {
		l, err := t.skipChar()
		if err != nil {
			return 0, err
		}
		ret += l
	}

	if len(t.data) == 0 {
		return 0, ErrorUnexpected.New(t.pos)
	}
	t.data = t.data[1:]
	return ret, nil
}

func (t *lexer) skipChar() (int, error) {
	switch c := t.data[0]; {
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRune(t.data)
		if r == utf8.RuneError {
			return 0, ErrorRune.New(t.pos)
		}
		t.data = t.data[size:]
		return 1, nil
	case c != '\\':
		t.data = t.data[1:]
		return 1, nil
	}

	// hard case: c is backslash
	if len(t.data) < 2 {
		return 0, ErrorUnexpected.New(t.pos)
	}
	c := t.data[1]
	t.data = t.data[2:]

	switch c {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
		return 2, nil
	case 'x', 'u', 'U':
		n := 0
		switch c {
		case 'x':
			n = 2
		case 'u':
			n = 4
		case 'U':
			n = 8
		}
		var v rune
		if len(t.data) < n {
			return 0, ErrorUnexpected.New(t.pos)
		}
		for j := 0; j < n; j++ {
			x, ok := t.unhex(t.data[j])
			if !ok {
				return 0, ErrorUnexpected.New(t.pos)
			}
			v = v<<4 | x
		}
		t.data = t.data[n:]
		if c == 'x' {
			return 2 + n, nil
		}
		if !utf8.ValidRune(v) {
			return 0, ErrorRune.New(t.pos)
		}
		return 2 + n, nil
	case '0', '1', '2', '3', '4', '5', '6', '7':
		v := rune(c) - '0'
		if len(t.data) < 2 {
			return 0, ErrorUnexpected.New(t.pos)
		}
		for j := 0; j < 2; j++ { // one digit already; two more
			x := rune(t.data[j]) - '0'
			if x < 0 || x > 7 {
				return 0, ErrorUnexpected.New(t.pos)
			}
			v = (v << 3) | x
		}
		t.data = t.data[2:]
		if v > 255 {
			return 0, ErrorUnexpected.New(t.pos)
		}
		return 4, nil
	case '\'':
		return 0, ErrorWrongQuote.New(t.pos)
	default:
		return 0, ErrorUnexpected.New(t.pos)
	}
}

func (t *lexer) unhex(b byte) (v rune, ok bool) {
	c := rune(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return
}
