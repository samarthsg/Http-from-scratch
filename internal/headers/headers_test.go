package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header + termination
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "localhost:42069", headers.Get("HOST"))
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid multiple headers + termination
	headers = NewHeaders()
	data = []byte(
		"Host: localhost:42069\r\n" +
			"User-Agent: Mozilla Firefox\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "Mozilla Firefox", headers.Get("User-Agent"))
	assert.Equal(t, "*/*", headers.Get("Accept"))
	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Existing header matches parsed header
	headers = NewHeaders()
	headers.Set("Host", "localhost")

	data = []byte(
		"Host: example.com\r\n" +
			"User-Agent: Firefox\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost, example.com", headers.Get("Host"))
	assert.Equal(t, "Firefox", headers.Get("User-Agent"))
	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Valid headers without termination
	headers = NewHeaders()
	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost", headers.Get("Host"))
	assert.Equal(t, "Firefox", headers.Get("User-Agent"))
	assert.False(t, done)
	assert.Equal(t, len(data), n)

	// Test: Existing headers
	headers = NewHeaders()
	headers.Set("Existing", "value")

	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "value", headers.Get("Existing"))
	assert.Equal(t, "localhost", headers.Get("Host"))
	assert.Equal(t, "Firefox", headers.Get("User-Agent"))
	assert.True(t, done)

	// Test: Extra whitespace around values
	headers = NewHeaders()
	data = []byte(
		"Host:       localhost:42069       \r\n" +
			"User-Agent:     Mozilla Firefox     \r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "Mozilla Firefox", headers.Get("User-Agent"))
	assert.True(t, done)

	// Test: Value containing multiple colons
	headers = NewHeaders()
	data = []byte("Host: localhost:42069:1234\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069:1234", headers.Get("Host"))
	assert.True(t, done)

	// Test: Value containing spaces
	headers = NewHeaders()
	data = []byte("User-Agent: Mozilla Firefox Browser\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "Mozilla Firefox Browser", headers.Get("User-Agent"))
	assert.True(t, done)

	// Test: Capitalized header names are case-insensitive
	headers = NewHeaders()
	data = []byte(
		"HoSt: localhost:42069\r\n" +
			"uSeR-AgEnT: Firefox\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "localhost:42069", headers.Get("HOST"))
	assert.Equal(t, "Firefox", headers.Get("User-Agent"))
	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Valid special characters in field name
	headers = NewHeaders()
	data = []byte("X!#$%&'*+-.^_`|~: value\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "value", headers.Get("X!#$%&'*+-.^_`|~"))
	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Invalid character in field name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Tab in field name
	headers = NewHeaders()
	data = []byte("Hos\t: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Leading whitespace before field name
	headers = NewHeaders()
	data = []byte("   Host: localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Whitespace before colon
	headers = NewHeaders()
	data = []byte("Host : localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Empty field name
	headers = NewHeaders()
	data = []byte(": localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Missing colon
	headers = NewHeaders()
	data = []byte("Host localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Empty field value
	headers = NewHeaders()
	data = []byte("Host:\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Incomplete final header
	headers = NewHeaders()
	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost", headers.Get("Host"))
	assert.Empty(t, headers.Get("User-Agent"))
	assert.False(t, done)
	assert.Equal(t, len("Host: localhost\r\n"), n)

	// Test: Empty input
	headers = NewHeaders()
	data = []byte{}

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers.Get("Host"))

	// Test: Only termination
	headers = NewHeaders()
	data = []byte("\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestHeaderSet(t *testing.T) {
	headers := NewHeaders()

	headers.Set("Content-Type", "Application/JSON")

	assert.Equal(t, "Application/JSON", headers.Get("Content-Type"))
	assert.Equal(t, "Application/JSON", headers.Get("content-type"))
	assert.Equal(t, "Application/JSON", headers.Get("CONTENT-TYPE"))
}
