package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

// HashString hashes the input string using SHA-256 and returns the hex representation.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// EncodeHash encodes the first numChars*2 characters of the hex hash into a base62 string of desired length.
// Returns an error if the hash is invalid or encoding fails.
func EncodeHash(hash string, numChars int) (string, error) {
	// Each byte is represented by two hex characters
	numHexChars := numChars * 2
	if len(hash) < numHexChars {
		return "", errors.New("hash length is insufficient for encoding")
	}

	// Convert hex string to bytes
	bytes, err := hex.DecodeString(hash[:numHexChars])
	if err != nil {
		return "", err
	}

	// Convert bytes to big.Int
	n := new(big.Int).SetBytes(bytes)

	// Encode the big.Int to base62 using EncodeBigInt
	encoded, err := EncodeBigInt(n)
	if err != nil {
		return "", err
	}

	// If encoded string is shorter than desired, pad with '0's
	if len(encoded) < numChars {
		encoded = strings.Repeat("0", numChars-len(encoded)) + encoded
	} else if len(encoded) > numChars {
		encoded = encoded[:numChars]
	}

	return encoded, nil
}
