package jajson

import (
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
	var tmp T
	size := int(unsafe.Sizeof(tmp)) * 8
	return getNumber[T](data, func(bytes []byte) (T, error) {
		//TODO replace ParseInt
		ret, err := strconv.ParseInt(string(bytes), 10, size)
		return T(ret), err
	}, path...)
}

func GetUInt[T uint | uint8 | uint16 | uint32 | uint64](data []byte, path ...string) (T, error) {
	var tmp T
	size := int(unsafe.Sizeof(tmp)) * 8
	return getNumber[T](data, func(bytes []byte) (T, error) {
		//TODO replace ParseUInt
		ret, err := strconv.ParseUint(string(bytes), 10, size)
		return T(ret), err
	}, path...)
}

func getNumber[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](data []byte, parse func([]byte) (T, error), path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return 0, err
	}
	if typ == Int {
		return parse(val)
	} else if typ == String {
		return parse(val[1 : len(val)-1])
	}
	return 0, ErrorWrongValueType
}

func GetFloat[T float32 | float64](data []byte, path ...string) (T, error) {
	typ, val, err := GetRawValue(data, path...)
	if err != nil {
		return 0, err
	}
	var tmp T
	//TODO replace ParseFloat
	if typ == Float {
		ret, err := strconv.ParseFloat(string(val), int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	} else if typ == String {
		ret, err := strconv.ParseFloat(string(val[1:len(val)-1]), int(unsafe.Sizeof(tmp))*8)
		return T(ret), err
	}
	return 0, ErrorWrongValueType
}
