package services

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBigInt(t *testing.T) {
	// Define test cases with corrected expected values
	tests := []struct {
		name        string
		input       string // Decimal string
		expected    string
		expectError bool
	}{
		{
			name:        "Encode BigInt 0",
			input:       "0",
			expected:    "0",
			expectError: false,
		},
		{
			name:        "Encode BigInt 1",
			input:       "1",
			expected:    "1",
			expectError: false,
		},
		{
			name:        "Encode BigInt 10",
			input:       "10",
			expected:    "A",
			expectError: false,
		},
		{
			name:        "Encode BigInt 61",
			input:       "61",
			expected:    "z",
			expectError: false,
		},
		{
			name:        "Encode BigInt 62",
			input:       "62",
			expected:    "10",
			expectError: false,
		},
		{
			name:        "Encode BigInt 12345",
			input:       "12345",
			expected:    "3D7",
			expectError: false,
		},
		{
			name:        "Encode BigInt Negative",
			input:       "-123456",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Encode BigInt Nil",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var n *big.Int
			if tt.input != "" {
				n = new(big.Int)
				_, ok := n.SetString(tt.input, 10)
				assert.True(t, ok, "SetString(%q) should succeed", tt.input)
			} else {
				n = nil
			}

			result, err := EncodeBigInt(n)
			if tt.expectError {
				assert.Error(t, err, "EncodeBigInt(%q) should return an error", tt.input)
			} else {
				assert.NoError(t, err, "EncodeBigInt(%q) should not return an error", tt.input)
				assert.Equal(t, tt.expected, result, "EncodeBigInt(%q) should be %q", tt.input, tt.expected)
			}
		})
	}
}
