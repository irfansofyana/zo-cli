package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSSE_BasicData(t *testing.T) {
	input := "data: hello\ndata: world\ndata: [DONE]\n"
	var chunks []string
	err := ReadSSE(strings.NewReader(input), func(data string) error {
		chunks = append(chunks, data)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"hello", "world"}, chunks)
}

func TestReadSSE_SkipsEmptyAndNonDataLines(t *testing.T) {
	input := "event: message\n\ndata: chunk1\n\n: comment\ndata: chunk2\ndata: [DONE]\n"
	var chunks []string
	err := ReadSSE(strings.NewReader(input), func(data string) error {
		chunks = append(chunks, data)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"chunk1", "chunk2"}, chunks)
}

func TestReadSSE_HandlerError(t *testing.T) {
	input := "data: first\ndata: second\n"
	err := ReadSSE(strings.NewReader(input), func(data string) error {
		return fmt.Errorf("stop")
	})
	assert.EqualError(t, err, "stop")
}

func TestReadSSE_EmptyStream(t *testing.T) {
	var chunks []string
	err := ReadSSE(strings.NewReader(""), func(data string) error {
		chunks = append(chunks, data)
		return nil
	})
	require.NoError(t, err)
	assert.Empty(t, chunks)
}

func TestReadSSE_NoDoneSentinel(t *testing.T) {
	input := "data: one\ndata: two\n"
	var chunks []string
	err := ReadSSE(strings.NewReader(input), func(data string) error {
		chunks = append(chunks, data)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, chunks)
}
