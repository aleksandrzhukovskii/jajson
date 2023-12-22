package jajson

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LexerSuite struct {
	suite.Suite
}

func TestLexer(t *testing.T) {
	suite.Run(t, new(LexerSuite))
}

func (t *LexerSuite) TestNextTokenOK() {
	testCase := `   [  {  ]  }  true  false 123 12 1 -123 -12 -1 123.5 -123,56  :  "abd" "☺" "\xFF" "\377" "\u1234" "\U00010111" "\U0001011111" "\a\b\f\n\r\t\v\\\"" "\a"  ,  `
	l := newLexer([]byte(testCase))
	check := []struct {
		pos     int
		bytePos int
		typ     LexemType
		value   []byte
	}{
		{pos: 3, bytePos: 3, typ: openBracket},
		{pos: 6, bytePos: 6, typ: openCurve},
		{pos: 9, bytePos: 9, typ: closeBracket},
		{pos: 12, bytePos: 12, typ: closeCurve},
		{pos: 15, bytePos: 15, typ: Bool, value: []byte("true")},
		{pos: 21, bytePos: 21, typ: Bool, value: []byte("false")},
		{pos: 27, bytePos: 27, typ: Int, value: []byte("123")},
		{pos: 31, bytePos: 31, typ: Int, value: []byte("12")},
		{pos: 34, bytePos: 34, typ: Int, value: []byte("1")},
		{pos: 36, bytePos: 36, typ: Int, value: []byte("-123")},
		{pos: 41, bytePos: 41, typ: Int, value: []byte("-12")},
		{pos: 45, bytePos: 45, typ: Int, value: []byte("-1")},
		{pos: 48, bytePos: 48, typ: Float, value: []byte("123.5")},
		{pos: 54, bytePos: 54, typ: Float, value: []byte("-123,56")},
		{pos: 63, bytePos: 63, typ: colon},
		{pos: 66, bytePos: 66, typ: String, value: []byte(`"abd"`)},
		{pos: 72, bytePos: 72, typ: String, value: []byte(`"☺"`)},
		{pos: 76, bytePos: 78, typ: String, value: []byte(`"\xFF"`)},
		{pos: 83, bytePos: 85, typ: String, value: []byte(`"\377"`)},
		{pos: 90, bytePos: 92, typ: String, value: []byte(`"\u1234"`)},
		{pos: 99, bytePos: 101, typ: String, value: []byte(`"\U00010111"`)},
		{pos: 112, bytePos: 114, typ: String, value: []byte(`"\U0001011111"`)},
		{pos: 127, bytePos: 129, typ: String, value: []byte(`"\a\b\f\n\r\t\v\\\""`)},
		{pos: 148, bytePos: 150, typ: String, value: []byte(`"\a"`)},
		{pos: 154, bytePos: 156, typ: comma},
	}
	for i := 0; i < len(check); i++ {
		lex, err := l.nextToken()
		t.NoError(err)
		t.Equal(check[i].typ, lex.typ)
		t.Equal(check[i].pos, lex.pos)
		if check[i].value != nil {
			t.Equal(check[i].value, lex.value)
		} else {
			t.Empty(lex.value)
		}
		if check[i].bytePos != 0 {
			t.Equal(check[i].bytePos, lex.bytePos)
		}
	}
	lex, err := l.nextToken()
	t.Equal(lexem{}, lex)
	t.EqualError(err, ErrorUnexpected.New(len([]rune(testCase))-1).Error())
}
