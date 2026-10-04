package headers

import (
	"bytes"
	"errors"
	"strings"
)

var seprator = []byte("\r\n")
var invalidFieldName = errors.New("Invalid field name")
var emptyFieldName = errors.New("Empty field name")
var invalidHeader = errors.New("Invalid header")
var emptyFieldValue = errors.New("Invalid field value")

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	key := strings.ToLower(name)
	if oldValue, exists := h.headers[key]; exists {
		h.headers[key] = oldValue + ", " + value
		return
	}
	h.headers[key] = value
}

func isValid(name string) bool {
	for _, ch := range name {
		found := false
		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

func extractInfo(header []byte) (string, string, error) {
	parts := bytes.SplitN(header, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", invalidHeader
	}
	fieldName := string(parts[0])
	fieldValue := string(bytes.TrimSpace(parts[1]))
	if len(fieldName) == 0 {
		return "", "", emptyFieldName
	}
	if !isValid(fieldName) {
		return "", "", invalidFieldName
	}
	if len(fieldValue) == 0 {
		return "", "", emptyFieldValue
	}
	return fieldName, fieldValue, nil
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
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
		h.Set(fieldName, fieldValue)
		read += len(seprator) + index
	}
	return read, done, nil
}
