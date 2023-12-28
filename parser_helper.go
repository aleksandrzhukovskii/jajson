package jajson

import (
	"bytes"
	"strconv"
)

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
		typ, data, err := checkLexeme(lex, before, closeCurve, Object, Error{})
		if typ != nothing {
			return typ, data, err
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

func checkLexeme(lex *lexer, before []byte, closeLexCheck, closeLexRet LexemeType, closeError Error) (LexemeType, []byte, error) {
	lxm, data, err := lex.nextToken()
	if err != nil {
		return Err, nil, err
	}
	if lxm.typ == closeLexCheck {
		if closeError.err != nil {
			return Err, nil, closeError.New(lxm.pos)
		}
		return closeLexRet, before[:len(before)-len(data)+1], nil
	} else if lxm.typ != comma {
		return Err, nil, ErrorUnexpectedLexeme.New(lxm.pos)
	}
	return nothing, before, nil
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
		typ, data, err := checkLexeme(lex, before, closeBracket, Array, Error{})
		if typ != nothing {
			return typ, data, err
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
		typ, _, err := checkLexeme(lex, nil, closeCurve, Err, ErrorWrongPath)
		if typ != nothing {
			return err
		}

		lxm, _, err := lex.nextToken()
		if err != nil {
			return err
		} else if lxm.typ != String {
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
