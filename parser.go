package jajson

import (
	"strconv"
	"unsafe"
)

// GetRawValue returns part of the original slice with value
func GetRawValue(data []byte, path ...string) (LexemType, []byte, error) {
	if len(data) == 0 {
		return 0, nil, ErrorEmptyJSON
	}
	lex := newLexer(data)
	if len(path) == 0 {
		return parseValue(lex)
	}
	return Err, nil, nil
}

func GetString(data []byte, path ...string) (string, error) {
	typ, val, err := GetRawValue(data, path...)
	if typ != String {
		return "", ErrorWrongValueType
	}
	if err != nil {
		return "", err
	}
	//TODO replace Unquote
	return strconv.Unquote(string(val))
}

func GetBool(data []byte, path ...string) (bool, error) {
	typ, val, err := GetRawValue(data, path...)
	if typ != Bool {
		return false, ErrorWrongValueType
	}
	if err != nil {
		return false, err
	}
	return val[0] == 't', nil
}

func GetInt[T int | int8 | int16 | int32 | int64](data []byte, path ...string) (T, error) {
	return getNumber[T](data, func(bytes []byte) (T, error) {
		var tmp T
		ret, err := strconv.ParseInt(string(bytes), 10, int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	}, path...)
}

func GetUInt[T uint | uint8 | uint16 | uint32 | uint64](data []byte, path ...string) (T, error) {
	return getNumber[T](data, func(bytes []byte) (T, error) {
		var tmp T
		ret, err := strconv.ParseUint(string(bytes), 10, int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	}, path...)
}

func getNumber[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](data []byte, parse func([]byte) (T, error), path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if typ != Int {
		return 0, ErrorWrongValueType
	}
	if err != nil {
		return 0, err
	}
	return parse(val)
}

func GetFloat[T float32 | float64](data []byte, path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if typ != Float {
		return 0, ErrorWrongValueType
	}
	if err != nil {
		return 0, err
	}
	var tmp T
	ret, err := strconv.ParseFloat(string(val), int(unsafe.Sizeof(tmp))*8)
	return T(ret), err
}

func parseValue(lex *lexer) (LexemType, []byte, error) {
	lexem, before, err := lex.nextToken()
	if err != nil {
		return Err, nil, err
	}
	if lexem.typ == String || lexem.typ == Bool || lexem.typ == Int || lexem.typ == Float {
		return lexem.typ, lexem.value, nil
	} else if lexem.typ == openCurve {
		return parseObject(lex, before)
	} else if lexem.typ == openBracket {
		return parseArray(lex, before)
	}
	return Err, nil, ErrorUnexpectedLexem.New(lexem.pos)
}

func parseObject(lex *lexer, before []byte) (LexemType, []byte, error) {
	//first bracket is skipped
	lexem, data, err := lex.nextToken()
	if err != nil {
		return Err, nil, err
	}
	if lexem.typ == closeCurve {
		return Object, before[:len(before)-len(data)+1], nil
	} else if lexem.typ != String {
		return Err, nil, ErrorUnexpectedLexem.New(lexem.pos)
	}

	if err := skipLexem(lex, colon); err != nil {
		return Err, nil, err
	}

	_, _, err = parseValue(lex)
	if err != nil {
		return Err, nil, err
	}
	return skipObjectFields(lex, before)
}

func skipObjectFields(lex *lexer, before []byte) (LexemType, []byte, error) {
	for {
		lexem, data, err := lex.nextToken()
		if err != nil {
			return Err, nil, err
		}
		if lexem.typ == closeCurve {
			return Object, before[:len(before)-len(data)+1], nil
		} else if lexem.typ != comma {
			return Err, nil, ErrorUnexpectedLexem.New(lexem.pos)
		}

		if err := skipLexem(lex, String); err != nil {
			return Err, nil, err
		}

		if err := skipLexem(lex, colon); err != nil {
			return Err, nil, err
		}

		_, _, err = parseValue(lex)
		if err != nil {
			return Err, nil, err
		}
	}
}

func skipLexem(lex *lexer, typ LexemType) error {
	lexem, _, err := lex.nextToken()
	if err != nil {
		return err
	}
	if lexem.typ != typ {
		return ErrorUnexpectedLexem.New(lexem.pos)
	}
	return nil
}

func parseArray(lex *lexer, before []byte) (LexemType, []byte, error) {
	//first bracket is skipped
	lexem, _, err := lex.lookup()
	if err != nil {
		return Err, nil, err
	}
	if lexem.typ == closeBracket {
		_, data, _ := lex.nextToken()
		return Array, before[:len(before)-len(data)+1], nil
	}
	if _, _, err = parseValue(lex); err != nil {
		return Err, nil, err
	}
	for {
		lexem, data, err := lex.nextToken()
		if err != nil {
			return Err, nil, err
		}
		if lexem.typ == closeBracket {
			return Array, before[:len(before)-len(data)+1], nil
		} else if lexem.typ != comma {
			return Err, nil, ErrorUnexpectedLexem.New(lexem.pos)
		}

		if _, _, err = parseValue(lex); err != nil {
			return Err, nil, err
		}
	}
}
