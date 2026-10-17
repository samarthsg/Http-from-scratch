package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
}

func (cr *chunkReader) Read(p []byte) (int, error) {
	if len(cr.data) == 0 {
		return 0, io.EOF
	}

	n := min(len(p), cr.numBytesPerRead, len(cr.data))
	copy(p, cr.data[:n])
	cr.data = cr.data[n:]

	return n, nil
}

func TestRequestFromReader(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		chunkSize int
		wantError bool
		method    string
		target    string
		version   string
		headers   map[string]string
		body      string
	}{
		{
			name:      "request with headers in small chunks",
			data:      "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			chunkSize: 3,
			method:    "GET",
			target:    "/",
			version:   "1.1",
			headers: map[string]string{
				"host":       "localhost:42069",
				"user-agent": "curl/7.81.0",
				"accept":     "*/*",
			},
			body: "",
		},
		{
			name:      "request with headers in one chunk",
			data:      "GET /hello HTTP/1.1\r\nHost: example.com\r\n\r\n",
			chunkSize: 1024,
			method:    "GET",
			target:    "/hello",
			version:   "1.1",
			headers: map[string]string{
				"host": "example.com",
			},
			body: "",
		},
		{
			name:      "request line only",
			data:      "GET / HTTP/1.1\r\n\r\n",
			chunkSize: 1024,
			method:    "GET",
			target:    "/",
			version:   "1.1",
			headers:   map[string]string{},
			body:      "",
		},
		{
			name:      "invalid request line",
			data:      "GET / HTTP/1.1 extra\r\n\r\n",
			chunkSize: 1024,
			wantError: true,
		},
		{
			name:      "unsupported HTTP version",
			data:      "GET / HTTP/2.0\r\n\r\n",
			chunkSize: 1024,
			wantError: true,
		},
		{
			name:      "invalid lowercase method",
			data:      "get / HTTP/1.1\r\n\r\n",
			chunkSize: 1024,
			wantError: true,
		},
		{
			name:      "incomplete request line",
			data:      "GET / HTTP/1.1",
			chunkSize: 3,
			wantError: true,
		},
		{
			name:      "incomplete headers",
			data:      "GET / HTTP/1.1\r\nHost: example.com\r\n",
			chunkSize: 3,
			wantError: true,
		},
		{
			name: "Standard Body",
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"hello world!\n",
			chunkSize: 3,
			method:    "POST",
			target:    "/submit",
			version:   "1.1",
			headers: map[string]string{
				"host":           "localhost:42069",
				"content-length": "13",
			},
			body: "hello world!\n",
		},
		{
			name: "Empty Body, 0 reported content length",
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 0\r\n" +
				"\r\n",
			chunkSize: 3,
			method:    "POST",
			target:    "/submit",
			version:   "1.1",
			headers: map[string]string{
				"host":           "localhost:42069",
				"content-length": "0",
			},
			body: "",
		},
		{
			name: "Empty Body, no reported content length",
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"\r\n",
			chunkSize: 3,
			method:    "POST",
			target:    "/submit",
			version:   "1.1",
			headers: map[string]string{
				"host": "localhost:42069",
			},
			body: "",
		},
		{
			name: "Body shorter than reported content length",
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 20\r\n" +
				"\r\n" +
				"partial content",
			chunkSize: 3,
			wantError: true,
		},
		{
			name: "No Content-Length but Body Exists",
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"\r\n" +
				"some body data",
			chunkSize: 3,
			method:    "POST",
			target:    "/submit",
			version:   "1.1",
			headers: map[string]string{
				"host": "localhost:42069",
			},
			body: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &chunkReader{
				data:            tt.data,
				numBytesPerRead: tt.chunkSize,
			}

			req, err := RequestFromReader(reader)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req)

			assert.Equal(t, tt.method, req.RequestLine.Method)
			assert.Equal(t, tt.target, req.RequestLine.RequestTarget)
			assert.Equal(t, tt.version, req.RequestLine.HttpVersion)

			for name, want := range tt.headers {
				got, exists := req.Headers.Get(name)
				assert.True(t, exists, "expected header %q to exist", name)
				assert.Equal(t, want, got)
			}

			assert.Equal(t, tt.body, string(req.Body))
		})
	}
}

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError error
		method    string
		target    string
		version   string
		wantN     int
	}{
		{
			name:    "valid request line",
			input:   "GET / HTTP/1.1\r\n",
			method:  "GET",
			target:  "/",
			version: "1.1",
			wantN:   len("GET / HTTP/1.1\r\n"),
		},
		{
			name:      "missing fields",
			input:     "GET /\r\n",
			wantError: InvalidRequest,
		},
		{
			name:      "invalid HTTP version",
			input:     "GET / HTTP/2.0\r\n",
			wantError: InvalidHttpRequest,
		},
		{
			name:      "invalid method",
			input:     "get / HTTP/1.1\r\n",
			wantError: InvalidMethod,
		},
		{
			name:  "incomplete line",
			input: "GET / HTTP/1.1",
			wantN: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, err, n := parseRequestLine([]byte(tt.input))

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantN, n)

			if tt.wantN == 0 {
				assert.Nil(t, line)
				return
			}

			require.NotNil(t, line)
			assert.Equal(t, tt.method, line.Method)
			assert.Equal(t, tt.target, line.RequestTarget)
			assert.Equal(t, tt.version, line.HttpVersion)
		})
	}
}
