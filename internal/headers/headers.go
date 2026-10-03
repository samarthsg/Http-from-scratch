package headers

import (
	"bytes"
	"errors"
)

var seprator = []byte("\r\n")
var invalidFieldName = errors.New("Invalid field name")
var emptyFieldName = errors.New("Empty field name")
var invalidHeader = errors.New("Invalid header")
var emptyFieldValue = errors.New("Invalid field value")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func extractInfo(header []byte) (string, string, error) {
	parts := bytes.SplitN(header, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", invalidHeader
	}
	fieldName := string(parts[0])
	fieldValue := string(bytes.TrimSpace(parts[1]))
	for _, ch := range fieldName {
		if ch == ' ' {
			return "", "", invalidFieldName
		}
	}
	if len(fieldName) == 0 {
		return "", "", emptyFieldName
	}
	if len(fieldValue) == 0 {
		return "", "", emptyFieldValue
	}
	return fieldName, fieldValue, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read, done := 0, false
	for {
		index := bytes.Index(data[read:], seprator)
		if index == -1 {
			break
		}
		if index == 0 {
			done = true
			read += len(seprator)
			break
		}
		fieldName, fieldValue, err := extractInfo(data[read : read+index])
		if err != nil {
			return read, false, err
		}
		h[fieldName] = fieldValue
		read += len(seprator) + index
	}
	return read, done, nil
}
