package request

import (
	"bytes"
	"errors"
	"io"
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
	bytesConsumed := 0
outer:
	for {
		switch r.State {
		case stateIntialized:
			rl, err, n := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			bytesConsumed += n
			r.State = stateDone
		case stateDone:
			break outer
		}
	}
	return bytesConsumed, nil
}

func (r *Request) isDone() bool {
	return r.State == stateDone
}

func parseRequestLine(data []byte) (*RequestLine, error, int) {
	var seprator = []byte("\r\n")
	request := bytes.Index(data, seprator)
	if request == -1 {
		return nil, nil, 0
	}
	firstLine := data[:request]
	totalParsed := request + len(seprator)
	metaData := bytes.Split(firstLine, []byte(" "))
	if len(metaData) != 3 {
		return nil, err, 0
	}
	httpParts := bytes.Split(metaData[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, InvalidHttpRequest, 0
	}
	ans := &RequestLine{
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(metaData[1]),
		Method:        string(metaData[0]),
	}
	if !isCapitalOnly(ans.Method) {
		return nil, InvalidMethod, 0
	}
	return ans, nil, totalParsed
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	request := newRequest()
	readToIndex := 0
	for !request.isDone() {
		if len(buff) == readToIndex {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}
		n, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			break
		}
		readToIndex += n
		readN, err := request.parse(buff[:readToIndex])
		copy(buff, buff[readN:readToIndex])
		readToIndex -= readN
	}
	return request, nil
}
