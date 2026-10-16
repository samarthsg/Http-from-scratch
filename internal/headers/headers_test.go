package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header + termination
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	value, exists := headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("HOST")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

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

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Mozilla Firefox", value)

	value, exists = headers.Get("Accept")
	assert.True(t, exists)
	assert.Equal(t, "*/*", value)

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

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost, example.com", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Firefox", value)

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

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Firefox", value)

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

	value, exists = headers.Get("Existing")
	assert.True(t, exists)
	assert.Equal(t, "value", value)

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Firefox", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Extra whitespace around values
	headers = NewHeaders()
	data = []byte(
		"Host:       localhost:42069       \r\n" +
			"User-Agent:     Mozilla Firefox     \r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Mozilla Firefox", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Value containing multiple colons
	headers = NewHeaders()
	data = []byte("Host: localhost:42069:1234\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069:1234", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Value containing spaces
	headers = NewHeaders()
	data = []byte("User-Agent: Mozilla Firefox Browser\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Mozilla Firefox Browser", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Capitalized header names are case-insensitive
	headers = NewHeaders()
	data = []byte(
		"HoSt: localhost:42069\r\n" +
			"uSeR-AgEnT: Firefox\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("host")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("HOST")
	assert.True(t, exists)
	assert.Equal(t, "localhost:42069", value)

	value, exists = headers.Get("User-Agent")
	assert.True(t, exists)
	assert.Equal(t, "Firefox", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Valid special characters in field name
	headers = NewHeaders()
	data = []byte("X!#$%&'*+-.^_`|~: value\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("X!#$%&'*+-.^_`|~")
	assert.True(t, exists)
	assert.Equal(t, "value", value)

	assert.True(t, done)
	assert.Equal(t, len(data), n)

	// Test: Invalid character in field name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Tab in field name
	headers = NewHeaders()
	data = []byte("Hos\t: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Leading whitespace before field name
	headers = NewHeaders()
	data = []byte("   Host: localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Whitespace before colon
	headers = NewHeaders()
	data = []byte("Host : localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Empty field name
	headers = NewHeaders()
	data = []byte(": localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Missing colon
	headers = NewHeaders()
	data = []byte("Host localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Empty field value
	headers = NewHeaders()
	data = []byte("Host:\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

	// Test: Incomplete final header
	headers = NewHeaders()
	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)

	value, exists = headers.Get("Host")
	assert.True(t, exists)
	assert.Equal(t, "localhost", value)

	value, exists = headers.Get("User-Agent")
	assert.False(t, exists)
	assert.Empty(t, value)

	assert.False(t, done)
	assert.Equal(t, len("Host: localhost\r\n"), n)

	// Test: Empty input
	headers = NewHeaders()
	data = []byte{}

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	value, exists = headers.Get("Host")
	assert.False(t, exists)
	assert.Empty(t, value)

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

	value, exists := headers.Get("Content-Type")
	assert.True(t, exists)
	assert.Equal(t, "Application/JSON", value)

	value, exists = headers.Get("content-type")
	assert.True(t, exists)
	assert.Equal(t, "Application/JSON", value)

	value, exists = headers.Get("CONTENT-TYPE")
	assert.True(t, exists)
	assert.Equal(t, "Application/JSON", value)
}

func TestHeaderGetMissing(t *testing.T) {
	headers := NewHeaders()

	value, exists := headers.Get("Content-Length")

	assert.False(t, exists)
	assert.Empty(t, value)
}
