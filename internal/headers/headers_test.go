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
	assert.Equal(t, "localhost:42069", headers["Host"])
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
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "Mozilla Firefox", headers["User-Agent"])
	assert.Equal(t, "*/*", headers["Accept"])
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
	assert.Equal(t, "localhost", headers["Host"])
	assert.Equal(t, "Firefox", headers["User-Agent"])
	assert.False(t, done)
	assert.Equal(t, len(data), n)

	// Test: Existing headers
	headers = NewHeaders()
	headers["Existing"] = "value"

	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox\r\n" +
			"\r\n",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "value", headers["Existing"])
	assert.Equal(t, "localhost", headers["Host"])
	assert.Equal(t, "Firefox", headers["User-Agent"])
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
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "Mozilla Firefox", headers["User-Agent"])
	assert.True(t, done)

	// Test: Value containing multiple colons
	headers = NewHeaders()
	data = []byte("Host: localhost:42069:1234\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069:1234", headers["Host"])
	assert.True(t, done)

	// Test: Value containing spaces
	headers = NewHeaders()
	data = []byte("User-Agent: Mozilla Firefox Browser\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "Mozilla Firefox Browser", headers["User-Agent"])
	assert.True(t, done)

	// Test: Leading whitespace before field name
	headers = NewHeaders()
	data = []byte("   Host: localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Whitespace before colon
	headers = NewHeaders()
	data = []byte("Host : localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Empty field name
	headers = NewHeaders()
	data = []byte(": localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Missing colon
	headers = NewHeaders()
	data = []byte("Host localhost\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Empty field value
	headers = NewHeaders()
	data = []byte("Host:\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Incomplete final header
	headers = NewHeaders()
	data = []byte(
		"Host: localhost\r\n" +
			"User-Agent: Firefox",
	)

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost", headers["Host"])
	assert.Empty(t, headers["User-Agent"])
	assert.False(t, done)
	assert.Equal(t, len("Host: localhost\r\n"), n)

	// Test: Empty input
	headers = NewHeaders()
	data = []byte{}

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Empty(t, headers)

	// Test: Only termination
	headers = NewHeaders()
	data = []byte("\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}
