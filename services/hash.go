package services

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// HashString hashes the input string using SHA-256
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// EncodeHash encodes the hash into a base62 string of desired length
func EncodeHash(hash string, numChars int) string {
	// Each byte is represented by two hex characters
	numHexChars := numChars * 2
	if len(hash) < numHexChars {
		numHexChars = len(hash)
	}

	// Convert hex string to bytes
	bytes, err := hex.DecodeString(hash[:numHexChars])
	if err != nil {
		// Handle error
		return ""
	}

	// Convert bytes to big.Int
	n := new(big.Int).SetBytes(bytes)

	// Encode the big.Int to base62
	return EncodeBigInt(n)
}
