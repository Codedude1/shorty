package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty String",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "Simple String",
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := HashString(tt.input)
			assert.Equal(t, tt.expected, result, "HashString(%q) should be %q", tt.input, tt.expected)
		})
	}
}

func TestEncodeHash(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		hash        string
		numChars    int
		expected    string
		expectError bool
	}{
		{
			name:     "Valid EncodeHash",
			hash:     "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			numChars: 6,
			// First 12 hex chars: "2cf24dba5fb0" => bytes: [0x2c, 0xf2, 0x4d, 0xba, 0x5f, 0xb0]
			// big.Int: 0x2cf24dba5fb0 = 4841617771328
			// base62 encoding of 4841617771328 is "E23GXW"
			expected:    "E23GXW",
			expectError: false,
		},
		{
			name:        "Hash Too Short",
			hash:        "abcdef",
			numChars:    4,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid Hex in Hash",
			hash:        "zzzzzzzzzzzz",
			numChars:    6,
			expected:    "",
			expectError: true,
		},
		{
			name:     "EncodeHash with Desired Length Longer than Encoding",
			hash:     "000000000000",
			numChars: 6,
			// First 12 hex chars: "000000000000" => bytes: [0x00, 0x00, 0x00, 0x00, 0x00, 0x00]
			// big.Int: 0
			// base62 encoding: "0"
			// Pad with '0's to make it 6 chars: "000000"
			expected:    "000000",
			expectError: false,
		},
		{
			name:     "EncodeHash with Zero Value",
			hash:     "000000000000",
			numChars: 1,
			// base62 encoding: "0"
			expected:    "0",
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncodeHash(tt.hash, tt.numChars)
			if tt.expectError {
				assert.Error(t, err, "EncodeHash(%q, %d) should return an error", tt.hash, tt.numChars)
			} else {
				assert.NoError(t, err, "EncodeHash(%q, %d) should not return an error", tt.hash, tt.numChars)
				assert.Equal(t, tt.expected, result, "EncodeHash(%q, %d) should be %q", tt.hash, tt.numChars, tt.expected)
			}
		})
	}
}

func TestEncodeHashConcurrency(t *testing.T) {
	// Define multiple hash strings and numChars
	hashes := []struct {
		hash     string
		numChars int
	}{
		{"2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", 6},
		{"ffffffffffff", 6},
		{"000000000000", 6},
	}

	done := make(chan bool)

	for _, h := range hashes {
		h := h // Capture range variable
		go func() {
			result, err := EncodeHash(h.hash, h.numChars)
			assert.NoError(t, err, "EncodeHash(%q, %d) should not return an error", h.hash, h.numChars)
			assert.NotEmpty(t, result, "EncodeHash(%q, %d) should return a non-empty string", h.hash, h.numChars)
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for range hashes {
		<-done
	}
}
