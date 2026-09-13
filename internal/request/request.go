package request

import (
	"errors"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

func isCapitalOnly(a string) bool {
	if len(a) == 0 {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] < 'A' || a[i] > 'Z' {
			return false
		}
	}
	return true
}

func parseRequestLine(part string) (*RequestLine, error) {
	request := strings.Split(part, " ")
	var err = errors.New("http: invalid request")
	if len(request) != 3 {
		return nil, err
	}
	version := request[2]
	version = version[len(version)-3:]
	ans := &RequestLine{
		HttpVersion:   version,
		RequestTarget: request[1],
		Method:        request[0],
	}
	var InvalidHttpRequest = errors.New("Incompatible Version found")
	var InvalidMethod = errors.New("Invalid method")
	if ans.HttpVersion != "1.1" {
		return nil, InvalidHttpRequest
	}
	if !isCapitalOnly(ans.Method) {
		return nil, InvalidMethod
	}
	return ans, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	allData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(allData), "\r\n")
	ans, parsererr := parseRequestLine(parts[0])
	if parsererr != nil {
		return nil, parsererr
	}
	req := &Request{RequestLine: *ans}
	return req, nil
}
