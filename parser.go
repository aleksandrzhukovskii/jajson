package jajson

import (
	"bytes"
	"strconv"
	"unsafe"
)

// GetRawValue returns part of the original slice with value
func GetRawValue(data []byte, path ...string) (LexemeType, []byte, error) {
	if len(data) == 0 {
		return 0, nil, ErrorEmptyJSON
	}
	lex := newLexer(data)
	if len(path) > 0 {
		if err := skipPath(lex, path); err != nil {
			return Err, nil, err
		}
	}
	return parseValue(lex)
}

func GetString(data []byte, path ...string) (string, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return "", err
	}
	if typ != String {
		return "", ErrorWrongValueType
	}
	//TODO replace Unquote
	return strconv.Unquote(string(val))
}

func GetBool(data []byte, path ...string) (bool, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return false, err
	}
	if typ != Bool {
		return false, ErrorWrongValueType
	}
	return val[0] == 't', nil
}

func GetInt[T int | int8 | int16 | int32 | int64](data []byte, path ...string) (T, error) {
	return getNumber[T](data, func(bytes []byte) (T, error) {
		var tmp T
		//TODO replace ParseInt
		ret, err := strconv.ParseInt(string(bytes), 10, int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	}, path...)
}

func GetUInt[T uint | uint8 | uint16 | uint32 | uint64](data []byte, path ...string) (T, error) {
	return getNumber[T](data, func(bytes []byte) (T, error) {
		var tmp T
		//TODO replace ParseUInt
		ret, err := strconv.ParseUint(string(bytes), 10, int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	}, path...)
}

func getNumber[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](data []byte, parse func([]byte) (T, error), path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return 0, err
	}
	if typ != Int {
		return 0, ErrorWrongValueType
	}
	return parse(val)
}

func GetFloat[T float32 | float64](data []byte, path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return 0, err
	}
	if typ != Float {
		return 0, ErrorWrongValueType
	}
	var tmp T
	//TODO replace ParseFloat
	ret, err := strconv.ParseFloat(string(val), int(unsafe.Sizeof(tmp))*8)
	return T(ret), err
}

func parseValue(lex *lexer) (LexemeType, []byte, error) {
	lxm, before, err := lex.nextToken()
	if err != nil {
		return Err, nil, err
	}
	if lxm.typ == String || lxm.typ == Bool || lxm.typ == Int || lxm.typ == Float {
		return lxm.typ, lxm.value, nil
	} else if lxm.typ == openCurve {
		return parseObject(lex, before)
	} else if lxm.typ == openBracket {
		return parseArray(lex, before)
	}
	return Err, nil, ErrorUnexpectedLexeme.New(lxm.pos)
}

func parseObject(lex *lexer, before []byte) (LexemeType, []byte, error) {
	//first bracket is skipped
	lxm, data, err := lex.nextToken()
	if err != nil {
		return Err, nil, err
	}
	if lxm.typ == closeCurve {
		return Object, before[:len(before)-len(data)+1], nil
	} else if lxm.typ != String {
		return Err, nil, ErrorUnexpectedLexeme.New(lxm.pos)
	}

	if err := skipLexeme(lex, colon); err != nil {
		return Err, nil, err
	}

	_, _, err = parseValue(lex)
	if err != nil {
		return Err, nil, err
	}
	return skipObjectFields(lex, before)
}

func skipObjectFields(lex *lexer, before []byte) (LexemeType, []byte, error) {
	for {
		lxm, data, err := lex.nextToken()
		if err != nil {
			return Err, nil, err
		}
		if lxm.typ == closeCurve {
			return Object, before[:len(before)-len(data)+1], nil
		} else if lxm.typ != comma {
			return Err, nil, ErrorUnexpectedLexeme.New(lxm.pos)
		}

		if err := skipLexeme(lex, String); err != nil {
			return Err, nil, err
		}

		if err := skipLexeme(lex, colon); err != nil {
			return Err, nil, err
		}

		_, _, err = parseValue(lex)
		if err != nil {
			return Err, nil, err
		}
	}
}

func skipLexeme(lex *lexer, typ LexemeType) error {
	lxm, _, err := lex.nextToken()
	if err != nil {
		return err
	}
	if lxm.typ != typ {
		return ErrorUnexpectedLexeme.New(lxm.pos)
	}
	return nil
}

func parseArray(lex *lexer, before []byte) (LexemeType, []byte, error) {
	//first bracket is skipped
	lxm, _, err := lex.lookup()
	if err != nil {
		return Err, nil, err
	}
	if lxm.typ == closeBracket {
		_, data, _ := lex.nextToken()
		return Array, before[:len(before)-len(data)+1], nil
	}
	if _, _, err = parseValue(lex); err != nil {
		return Err, nil, err
	}
	for {
		lxm, data, err := lex.nextToken()
		if err != nil {
			return Err, nil, err
		}
		if lxm.typ == closeBracket {
			return Array, before[:len(before)-len(data)+1], nil
		} else if lxm.typ != comma {
			return Err, nil, ErrorUnexpectedLexeme.New(lxm.pos)
		}

		if _, _, err = parseValue(lex); err != nil {
			return Err, nil, err
		}
	}
}

func skipPath(lex *lexer, path []string) error {
	for i := 0; i < len(path); i++ {
		if err := skipPathPart(lex, []byte(path[i])); err != nil {
			return err
		}
	}
	return nil
}

func skipPathPart(lex *lexer, path []byte) error {
	if err := skipLexeme(lex, openCurve); err != nil {
		return err
	}

	lxm, _, err := lex.nextToken()
	if err != nil {
		return err
	}
	if lxm.typ != String {
		return ErrorUnexpectedLexeme.New(lxm.pos)
	}
	//TODO replace Unquote
	field, err := strconv.Unquote(string(lxm.value))
	if err != nil {
		return err
	}

	if err := skipLexeme(lex, colon); err != nil {
		return err
	}

	if bytes.Compare(path, []byte(field)) == 0 {
		return nil
	}

	if _, _, err := parseValue(lex); err != nil {
		return err
	}

	return skipPathPartFields(lex, path)
}

func skipPathPartFields(lex *lexer, path []byte) error {
	for {
		lxm, _, err := lex.nextToken()
		if err != nil {
			return err
		} else if lxm.typ == closeCurve {
			return ErrorWrongPath.New(lxm.pos)
		} else if lxm.typ != comma {
			return ErrorUnexpectedLexeme.New(lxm.pos)
		}

		lxm, _, err = lex.nextToken()
		if err != nil {
			return err
		}
		if lxm.typ != String {
			return ErrorUnexpectedLexeme.New(lxm.pos)
		}
		//TODO replace Unquote
		field, err := strconv.Unquote(string(lxm.value))
		if err != nil {
			return err
		}

		if err := skipLexeme(lex, colon); err != nil {
			return err
		}

		if bytes.Compare(path, []byte(field)) == 0 {
			return nil
		}

		if _, _, err := parseValue(lex); err != nil {
			return err
		}
	}
}
