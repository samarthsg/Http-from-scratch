package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

type parserState int

const (
	stateIntialized parserState = iota
	stateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       parserState
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

var allDone = errors.New("Error: trying to read data in a done state")
var err = errors.New("http: invalid request")
var InvalidHttpRequest = errors.New("Incompatible Version found")
var InvalidMethod = errors.New("Invalid method")

func newRequest() *Request {
	return &Request{
		State: stateIntialized,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == stateDone {
		return 0, allDone
	}
	r.State = stateDone
	return len(data), nil
}

func parseRequestLine(data []byte) (*RequestLine, error, int) {
	request := strings.Split(part, " ")
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
	if ans.HttpVersion != "1.1" {
		return nil, InvalidHttpRequest
	}
	if !isCapitalOnly(ans.Method) {
		return nil, InvalidMethod
	}
	return ans, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	request := newRequest()
	readToIndex := 0
	for {
		read, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			request.State = 1
			break
		}
		bytesRead, err := request.parse(buff)
		if err == allDone {
			fmt.Println("Parsed already : ", err)
		}

	}
	return request, nil
}
