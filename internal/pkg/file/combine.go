package file

import (
	"bytes"
	"errors"
)

var ErrEmptyParts = errors.New("empty parts")

func Combine(parts [][]byte) ([]byte, error) {
	if len(parts) == 0 {
		return nil, ErrEmptyParts
	}

	var result bytes.Buffer

	for i := range parts {
		_, err := result.Write(parts[i])
		if err != nil {
			return nil, err
		}
	}

	return result.Bytes(), nil

}
