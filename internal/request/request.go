package request

import (
	"bytes"
	"errors"
	"http-from-scratch/internal/headers"
	"io"
	"strconv"
)

type parserState int

const (
	stateIntialized parserState = iota
	stateHeaders
	stateBody
	stateDone
	stateError
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       parserState
	Headers     *headers.Headers
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		State:   stateIntialized,
		Headers: headers.NewHeaders(),
	}
}

var allDone = errors.New("Error: trying to read data in a done state")
var InvalidRequest = errors.New("http: invalid request")
var InvalidHttpRequest = errors.New("Incompatible Version found")
var InvalidMethod = errors.New("Invalid method")
var bufferOverFlow = errors.New("Buffer Overflow")
var contentLengthMismatch = errors.New("Content-Length mistmatch")

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

func getInt(headers *headers.Headers, name string, defaultValue int) (int, error) {
	strValue, exist := headers.Get(name)
	if !exist {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue, err
	}
	return value, nil
}

func (r *Request) parse(data []byte) (int, error) {
	bytesConsumed := 0
outer:
	for {
		currentData := data[bytesConsumed:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.State {
		case stateIntialized:
			rl, err, n := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			bytesConsumed += n
			r.State = stateHeaders
		case stateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.State = stateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			bytesConsumed += n
			if done {
				_, exists := r.Headers.Get("Content-Length")
				if !exists {
					r.State = stateDone
				} else {
					r.State = stateBody
				}
			}
		case stateBody:
			length, err := getInt(r.Headers, "Content-Length", 0)
			if err != nil {
				r.State = stateError
				return 0, err
			}
			if length < 0 {
				r.State = stateError
				return 0, contentLengthMismatch
			}
			if length == 0 {
				r.State = stateDone
				break outer
			}
			remaining := min(length-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remaining]...)
			bytesConsumed += remaining
			if len(r.Body) == length {
				r.State = stateDone
			}
		case stateDone, stateError:
			break outer
		}
	}
	return bytesConsumed, nil
}

func (r *Request) isDone() bool {
	return r.State == stateDone || r.State == stateError
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
		return nil, InvalidRequest, 0
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
	buff := make([]byte, 8196)
	request := newRequest()
	readToIndex := 0
	for !request.isDone() {
		if len(buff) == readToIndex {
			return nil, bufferOverFlow
		}
		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			return nil, err
		}
		readToIndex += n
		readN, err := request.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[readN:readToIndex])
		readToIndex -= readN
	}
	return request, nil
}
