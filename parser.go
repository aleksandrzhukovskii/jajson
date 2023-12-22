package jajson

import "errors"

// GetRawValue returns part of original slice with the value
func GetRawValue(data []byte, path ...string) (LexemType, []byte, error) {
	if len(data) == 0 {
		return 0, nil, ErrorEmptyJSON
	}

	return String, nil, nil
}

func skipPath(lexer lexer, path []string) error {
	return errors.New("Not implemented")
}
