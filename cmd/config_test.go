package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"zo_sk_abc123xyz", "***********3xyz"}, // Mask all but last 4
		{"ab", "****"},                          // Short key
		{"", "****"},                            // Empty
		{"abcd", "****"},                        // Exactly 4
		{"abcde", "*bcde"},                      // 5 chars
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, maskKey(tt.input))
		})
	}
}
